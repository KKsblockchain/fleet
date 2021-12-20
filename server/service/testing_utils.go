package service

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/WatchBeam/clock"
	eeservice "github.com/fleetdm/fleet/v4/ee/server/service"
	"github.com/fleetdm/fleet/v4/server/config"
	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/fleetdm/fleet/v4/server/logging"
	"github.com/fleetdm/fleet/v4/server/ptr"
	"github.com/fleetdm/fleet/v4/server/service/async"
	"github.com/fleetdm/fleet/v4/server/sso"
	kitlog "github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/throttled/throttled/v2/store/memstore"
)

func newTestService(ds fleet.Datastore, rs fleet.QueryResultStore, lq fleet.LiveQueryStore, opts ...TestServerOpts) fleet.Service {
	return newTestServiceWithConfig(ds, config.TestConfig(), rs, lq, opts...)
}

func newTestServiceWithConfig(ds fleet.Datastore, fleetConfig config.FleetConfig, rs fleet.QueryResultStore, lq fleet.LiveQueryStore, opts ...TestServerOpts) fleet.Service {
	mailer := &mockMailService{SendEmailFn: func(e fleet.Email) error { return nil }}
	license := &fleet.LicenseInfo{Tier: fleet.TierFree}
	writer, _ := logging.NewFilesystemLogWriter(
		fleetConfig.Filesystem.StatusLogFile,
		kitlog.NewNopLogger(),
		fleetConfig.Filesystem.EnableLogRotation,
		fleetConfig.Filesystem.EnableLogCompression,
	)
	// See #1776
	//if err != nil {
	//	panic(err)
	//}
	osqlogger := &logging.OsqueryLogger{Status: writer, Result: writer}
	logger := kitlog.NewNopLogger()

	var ssoStore sso.SessionStore
	if len(opts) > 0 {
		if opts[0].Logger != nil {
			logger = opts[0].Logger
		}
		if opts[0].License != nil {
			license = opts[0].License
		}
		if opts[0].Pool != nil {
			ssoStore = sso.NewSessionStore(opts[0].Pool)
		}
	}
	task := &async.Task{
		Datastore:    ds,
		AsyncEnabled: false,
	}
	svc, err := NewService(ds, task, rs, logger, osqlogger, fleetConfig, mailer, clock.C, ssoStore, lq, ds, *license, NewMemFailingPolicySet())
	if err != nil {
		panic(err)
	}
	if license.IsPremium() {
		svc, err = eeservice.NewService(svc, ds, kitlog.NewNopLogger(), fleetConfig, mailer, clock.C, license)
		if err != nil {
			panic(err)
		}
	}
	return svc
}

func newTestServiceWithClock(ds fleet.Datastore, rs fleet.QueryResultStore, lq fleet.LiveQueryStore, c clock.Clock) fleet.Service {
	mailer := &mockMailService{SendEmailFn: func(e fleet.Email) error { return nil }}
	license := fleet.LicenseInfo{Tier: fleet.TierFree}
	testConfig := config.TestConfig()
	writer, err := logging.NewFilesystemLogWriter(
		testConfig.Filesystem.StatusLogFile,
		kitlog.NewNopLogger(),
		testConfig.Filesystem.EnableLogRotation,
		testConfig.Filesystem.EnableLogCompression,
	)
	if err != nil {
		panic(err)
	}
	osqlogger := &logging.OsqueryLogger{Status: writer, Result: writer}
	task := &async.Task{
		Datastore:    ds,
		AsyncEnabled: false,
	}
	svc, err := NewService(ds, task, rs, kitlog.NewNopLogger(), osqlogger, testConfig, mailer, c, nil, lq, ds, license, NewMemFailingPolicySet())
	if err != nil {
		panic(err)
	}
	return svc
}

func createTestUsers(t *testing.T, ds fleet.Datastore) map[string]fleet.User {
	users := make(map[string]fleet.User)
	for _, u := range testUsers {
		role := fleet.RoleObserver
		if strings.Contains(u.Email, "admin") {
			role = fleet.RoleAdmin
		}
		user := &fleet.User{
			Name:       "Test Name " + u.Email,
			Email:      u.Email,
			GlobalRole: &role,
		}
		err := user.SetPassword(u.PlaintextPassword, 10, 10)
		require.Nil(t, err)
		user, err = ds.NewUser(context.Background(), user)
		require.Nil(t, err)
		users[user.Email] = *user
	}
	return users
}

var testUsers = map[string]struct {
	Email             string
	PlaintextPassword string
	GlobalRole        *string
}{
	"admin1": {
		PlaintextPassword: "foobarbaz1234!",
		Email:             "admin1@example.com",
		GlobalRole:        ptr.String(fleet.RoleAdmin),
	},
	"user1": {
		PlaintextPassword: "foobarbaz1234!",
		Email:             "user1@example.com",
		GlobalRole:        ptr.String(fleet.RoleMaintainer),
	},
	"user2": {
		PlaintextPassword: "bazfoo1234!",
		Email:             "user2@example.com",
		GlobalRole:        ptr.String(fleet.RoleObserver),
	},
}

type mockMailService struct {
	SendEmailFn func(e fleet.Email) error
	Invoked     bool
}

func (svc *mockMailService) SendEmail(e fleet.Email) error {
	svc.Invoked = true
	return svc.SendEmailFn(e)
}

type TestServerOpts struct {
	Logger              kitlog.Logger
	License             *fleet.LicenseInfo
	SkipCreateTestUsers bool
	Rs                  fleet.QueryResultStore
	Lq                  fleet.LiveQueryStore
	Pool                fleet.RedisPool
}

func RunServerForTestsWithDS(t *testing.T, ds fleet.Datastore, opts ...TestServerOpts) (map[string]fleet.User, *httptest.Server) {
	var rs fleet.QueryResultStore
	if len(opts) > 0 && opts[0].Rs != nil {
		rs = opts[0].Rs
	}
	var lq fleet.LiveQueryStore
	if len(opts) > 0 && opts[0].Lq != nil {
		lq = opts[0].Lq
	}
	svc := newTestService(ds, rs, lq, opts...)
	users := map[string]fleet.User{}
	if len(opts) == 0 || (len(opts) > 0 && !opts[0].SkipCreateTestUsers) {
		users = createTestUsers(t, ds)
	}
	logger := kitlog.NewLogfmtLogger(os.Stdout)
	if len(opts) > 0 && opts[0].Logger != nil {
		logger = opts[0].Logger
	}

	limitStore, _ := memstore.New(0)
	r := MakeHandler(svc, config.FleetConfig{}, logger, limitStore)
	server := httptest.NewServer(r)
	t.Cleanup(func() {
		server.Close()
	})
	return users, server
}

func testKinesisPluginConfig() config.FleetConfig {
	c := config.TestConfig()
	c.Filesystem = config.FilesystemConfig{}
	c.Osquery.ResultLogPlugin = "kinesis"
	c.Osquery.StatusLogPlugin = "kinesis"
	c.Kinesis = config.KinesisConfig{
		Region:           "us-east-1",
		AccessKeyID:      "foo",
		SecretAccessKey:  "bar",
		StsAssumeRoleArn: "baz",
		StatusStream:     "test-status-stream",
		ResultStream:     "test-result-stream",
	}
	return c
}

func testFirehosePluginConfig() config.FleetConfig {
	c := config.TestConfig()
	c.Filesystem = config.FilesystemConfig{}
	c.Osquery.ResultLogPlugin = "firehose"
	c.Osquery.StatusLogPlugin = "firehose"
	c.Firehose = config.FirehoseConfig{
		Region:           "us-east-1",
		AccessKeyID:      "foo",
		SecretAccessKey:  "bar",
		StsAssumeRoleArn: "baz",
		StatusStream:     "test-status-firehose",
		ResultStream:     "test-result-firehose",
	}
	return c
}

func testLambdaPluginConfig() config.FleetConfig {
	c := config.TestConfig()
	c.Filesystem = config.FilesystemConfig{}
	c.Osquery.ResultLogPlugin = "lambda"
	c.Osquery.StatusLogPlugin = "lambda"
	c.Lambda = config.LambdaConfig{
		Region:           "us-east-1",
		AccessKeyID:      "foo",
		SecretAccessKey:  "bar",
		StsAssumeRoleArn: "baz",
		ResultFunction:   "result-func",
		StatusFunction:   "status-func",
	}
	return c
}

func testPubSubPluginConfig() config.FleetConfig {
	c := config.TestConfig()
	c.Filesystem = config.FilesystemConfig{}
	c.Osquery.ResultLogPlugin = "pubsub"
	c.Osquery.StatusLogPlugin = "pubsub"
	c.PubSub = config.PubSubConfig{
		Project:       "test",
		StatusTopic:   "status-topic",
		ResultTopic:   "result-topic",
		AddAttributes: false,
	}
	return c
}

func testStdoutPluginConfig() config.FleetConfig {
	c := config.TestConfig()
	c.Filesystem = config.FilesystemConfig{}
	c.Osquery.ResultLogPlugin = "stdout"
	c.Osquery.StatusLogPlugin = "stdout"
	return c
}

func assertBodyContains(t *testing.T, resp *http.Response, expected string) {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	require.Nil(t, err)
	bodyString := string(bodyBytes)
	assert.Contains(t, bodyString, expected)
}

func getJSON(r *http.Response, target interface{}) error {
	return json.NewDecoder(r.Body).Decode(target)
}

func assertErrorCodeAndMessage(t *testing.T, resp *http.Response, code int, message string) {
	err := &fleet.Error{}
	require.Nil(t, getJSON(resp, err))
	assert.Equal(t, code, err.Code)
	assert.Equal(t, message, err.Message)
}

type memFailingPolicySet struct {
	mMu sync.RWMutex
	m   map[uint][]PolicySetHost
}

var _ FailingPolicySet = (*memFailingPolicySet)(nil)

func NewMemFailingPolicySet() *memFailingPolicySet {
	return &memFailingPolicySet{
		m: make(map[uint][]PolicySetHost),
	}
}

// AddFailingPoliciesForHost adds the given host to the policy sets.
func (m *memFailingPolicySet) AddHost(policyID uint, host PolicySetHost) error {
	m.mMu.Lock()
	defer m.mMu.Unlock()

	m.m[policyID] = append(m.m[policyID], host)
	return nil
}

// ListHosts returns the list of hosts present in the policy set.
func (m *memFailingPolicySet) ListHosts(policyID uint) ([]PolicySetHost, error) {
	m.mMu.RLock()
	defer m.mMu.RUnlock()

	hosts := make([]PolicySetHost, len(m.m[policyID]))
	for i := range m.m[policyID] {
		hosts[i] = m.m[policyID][i]
	}
	return hosts, nil
}

// RemoveHosts removes the hosts from the policy set.
func (m *memFailingPolicySet) RemoveHosts(policyID uint, hosts []PolicySetHost) error {
	m.mMu.Lock()
	defer m.mMu.Unlock()

	if _, ok := m.m[policyID]; !ok {
		return nil
	}
	hostsSet := make(map[uint]struct{})
	for _, host := range hosts {
		hostsSet[host.ID] = struct{}{}
	}
	n := 0
	for _, host := range m.m[policyID] {
		if _, ok := hostsSet[host.ID]; !ok {
			m.m[policyID][n] = host
			n++
		}
	}
	m.m[policyID] = m.m[policyID][:n]
	return nil
}

// RemoveSet removes a policy set.
func (m *memFailingPolicySet) RemoveSet(policyID uint) error {
	m.mMu.Lock()
	defer m.mMu.Unlock()

	delete(m.m, policyID)
	return nil
}

// ListSets lists all the policy sets.
func (m *memFailingPolicySet) ListSets() ([]uint, error) {
	m.mMu.RLock()
	defer m.mMu.RUnlock()

	var policyIDs []uint
	for policyID := range m.m {
		policyIDs = append(policyIDs, policyID)
	}
	return policyIDs, nil
}

func RunFailingPolicySetTests(t *testing.T, r FailingPolicySet) {
	policyID1 := uint(1)

	// Test listing policy sets with no sets.
	policyIDs, err := r.ListSets()
	require.NoError(t, err)
	require.Empty(t, policyIDs)

	// Test listing if the policy set doesn't exist.
	hostIDs, err := r.ListHosts(policyID1)
	require.NoError(t, err)
	require.Empty(t, hostIDs)

	// Test removing hosts if the set doesn't exist.
	hostx := PolicySetHost{
		ID:       uint(999),
		Hostname: "hostx.example",
	}
	err = r.RemoveHosts(policyID1, []PolicySetHost{hostx})
	require.NoError(t, err)

	// Test adding a new host to a policy set.
	host2 := PolicySetHost{
		ID:       uint(2),
		Hostname: "host2.example",
	}
	err = r.AddHost(policyID1, host2)
	require.NoError(t, err)

	// Test listing the created policy set.
	policyIDs, err = r.ListSets()
	require.NoError(t, err)
	require.Len(t, policyIDs, 1)
	require.Equal(t, policyID1, policyIDs[0])

	// Test listing the policy set.
	hostIDs, err = r.ListHosts(policyID1)
	require.NoError(t, err)
	require.Len(t, hostIDs, 1)
	require.Equal(t, host2, hostIDs[0])

	// Test adding a second host to the policy set.
	host3 := PolicySetHost{
		ID:       uint(3),
		Hostname: "host3.example",
	}
	err = r.AddHost(policyID1, host3)
	require.NoError(t, err)

	// Test adding a shared host on a different policy set.
	policyID2 := uint(2)
	err = r.AddHost(policyID2, host2)
	require.NoError(t, err)

	// Test listing the newly created policy set.
	policyIDs, err = r.ListSets()
	require.NoError(t, err)
	require.Len(t, policyIDs, 2)
	sort.Slice(policyIDs, func(i, j int) bool {
		return policyIDs[i] < policyIDs[j]
	})
	require.Equal(t, policyID1, policyIDs[0])
	require.Equal(t, policyID2, policyIDs[1])

	// Test listing of first and second policy set.
	hostIDs, err = r.ListHosts(policyID1)
	require.NoError(t, err)
	require.Len(t, hostIDs, 2)
	sort.Slice(hostIDs, func(i, j int) bool {
		return hostIDs[i].ID < hostIDs[j].ID
	})
	require.Equal(t, host2, hostIDs[0])
	require.Equal(t, host3, hostIDs[1])
	hostIDs, err = r.ListHosts(policyID2)
	require.NoError(t, err)
	require.Len(t, hostIDs, 1)
	require.Equal(t, host2, hostIDs[0])

	// Test removing all hosts from the first policy set.
	err = r.RemoveHosts(policyID1, []PolicySetHost{host2, host3})
	require.NoError(t, err)

	// Test listing of first and second policy set (after some removal).
	hostIDs, err = r.ListHosts(policyID1)
	require.NoError(t, err)
	require.Empty(t, hostIDs)
	hostIDs, err = r.ListHosts(policyID2)
	require.NoError(t, err)
	require.Len(t, hostIDs, 1)
	require.Equal(t, host2, hostIDs[0])

	// Test removal of a host on the second policy set.
	err = r.RemoveHosts(policyID2, []PolicySetHost{host2})
	require.NoError(t, err)

	// Test listing of the second policy set.
	hostIDs, err = r.ListHosts(policyID2)
	require.NoError(t, err)
	require.Empty(t, hostIDs)

	// Add another host to the first policy set (after removal of the other two).
	host4 := PolicySetHost{
		ID:       uint(4),
		Hostname: "host4.example",
	}
	err = r.AddHost(policyID1, host4)
	require.NoError(t, err)

	// Test listing of the first policy set.
	hostIDs, err = r.ListHosts(policyID1)
	require.NoError(t, err)
	require.Len(t, hostIDs, 1)
	require.Equal(t, host4, hostIDs[0])

	// Test removal of the first policy set.
	err = r.RemoveSet(policyID1)
	require.NoError(t, err)

	// Test listing of a removed policy set.
	hostIDs, err = r.ListHosts(policyID1)
	require.NoError(t, err)
	require.Empty(t, hostIDs)

	// Test a second removal of the first policy set.
	err = r.RemoveSet(policyID1)
	require.NoError(t, err)

	// Test removing hosts if the set doesn't exist anymore.
	err = r.RemoveHosts(policyID1, []PolicySetHost{host4})
	require.NoError(t, err)

	// Test removal of an unexisting policy set.
	policyIDX := uint(999)
	err = r.RemoveSet(policyIDX)
	require.NoError(t, err)

	// Test listing policy still returns the policyID2 set.
	policyIDs, err = r.ListSets()
	require.NoError(t, err)
	require.Len(t, policyIDs, 1)
	require.Equal(t, policyID2, policyIDs[0])

	// Now remove the remaining set.
	err = r.RemoveSet(policyID2)
	require.NoError(t, err)

	// And now it should be empty.
	policyIDs, err = r.ListSets()
	require.NoError(t, err)
	require.Empty(t, policyIDs)
}
