// Used in AddPackageModal.tsx and EditSoftwareModal.tsx
import React, { useContext, useState } from "react";

import { NotificationContext } from "context/notification";
import { getFileDetails } from "utilities/file/fileUtils";
import getInstallScript from "utilities/software_install_scripts";

import Button from "components/buttons/Button";
import Checkbox from "components/forms/fields/Checkbox";
import { FileUploader } from "components/FileUploader/FileUploader";
import Spinner from "components/Spinner";
import TooltipWrapper from "components/TooltipWrapper";

import AddPackageAdvancedOptions from "../AddPackageAdvancedOptions";

import { generateFormValidation } from "./helpers";

export const baseClass = "add-package-form";

const UploadingSoftware = () => {
  return (
    <div className={`${baseClass}__uploading-message`}>
      <Spinner centered={false} />
      <p>Adding software. This may take a few minutes to finish.</p>
    </div>
  );
};

export interface IAddPackageFormData {
  software: File | null;
  installScript: string;
  preInstallQuery?: string;
  postInstallScript?: string;
  uninstallScript?: string;
  selfService: boolean;
}

export interface IFormValidation {
  isValid: boolean;
  software: { isValid: boolean };
  preInstallQuery?: { isValid: boolean; message?: string };
  postInstallScript?: { isValid: boolean; message?: string };
  uninstallScript?: { isValid: boolean; message?: string };
  selfService?: { isValid: boolean };
}

interface IAddPackageFormProps {
  isUploading: boolean;
  onCancel: () => void;
  onSubmit: (formData: IAddPackageFormData) => void;
  isEditingSoftware?: boolean;
  defaultSoftware?: any; // TODO
  defaultInstallScript?: string;
  defaultPreInstallQuery?: string;
  defaultPostInstallScript?: string;
  defaultUninstallScript?: string;
  defaultSelfService?: boolean;
}

const ACCEPTED_EXTENSIONS = ".pkg,.msi,.exe,.deb";

const AddPackageForm = ({
  isUploading,
  onCancel,
  onSubmit,
  isEditingSoftware = false,
  defaultSoftware,
  defaultInstallScript,
  defaultPreInstallQuery,
  defaultPostInstallScript,
  defaultUninstallScript,
  defaultSelfService,
}: IAddPackageFormProps) => {
  const { renderFlash } = useContext(NotificationContext);

  const [formData, setFormData] = useState<IAddPackageFormData>({
    software: defaultSoftware || null,
    installScript: defaultInstallScript || "",
    preInstallQuery: defaultPreInstallQuery || undefined,
    postInstallScript: defaultPostInstallScript || undefined,
    uninstallScript: defaultUninstallScript || undefined,
    selfService: defaultSelfService || false,
  });
  const [formValidation, setFormValidation] = useState<IFormValidation>({
    isValid: false,
    software: { isValid: false },
  });

  const onFileUpload = (files: FileList | null) => {
    if (files && files.length > 0) {
      const file = files[0];

      let installScript: string;
      try {
        installScript = getInstallScript(file.name);
      } catch (e) {
        renderFlash("error", `${e}`);
        return;
      }

      const newData = {
        ...formData,
        software: file,
        installScript,
      };
      setFormData(newData);
      setFormValidation(generateFormValidation(newData));
    }
  };

  const onFormSubmit = (evt: React.FormEvent<HTMLFormElement>) => {
    evt.preventDefault();
    onSubmit(formData);
  };

  const onChangeInstallScript = (value: string) => {
    setFormData({ ...formData, installScript: value });
  };

  const onChangePreInstallQuery = (value?: string) => {
    const newData = { ...formData, preInstallQuery: value };
    setFormData(newData);
    setFormValidation(generateFormValidation(newData));
  };

  const onChangePostInstallScript = (value?: string) => {
    const newData = { ...formData, postInstallScript: value };
    setFormData(newData);
    setFormValidation(generateFormValidation(newData));
  };

  const onToggleSelfServiceCheckbox = (value: boolean) => {
    const newData = { ...formData, selfService: value };
    setFormData(newData);
    setFormValidation(generateFormValidation(newData));
  };

  const isSubmitDisabled = !formValidation.isValid;

  return (
    <div className={baseClass}>
      {isUploading ? (
        <UploadingSoftware />
      ) : (
        <form className={`${baseClass}__form`} onSubmit={onFormSubmit}>
          <FileUploader
            canEdit={isEditingSoftware}
            graphicName={"file-pkg"}
            accept={ACCEPTED_EXTENSIONS}
            message=".pkg, .msi, .exe, or .deb"
            onFileUpload={onFileUpload}
            buttonMessage="Choose file"
            buttonType="link"
            className={`${baseClass}__file-uploader`}
            fileDetails={
              formData.software ? getFileDetails(formData.software) : undefined
            }
          />
          <Checkbox
            value={formData.selfService}
            onChange={onToggleSelfServiceCheckbox}
          >
            <TooltipWrapper
              tipContent={
                <>
                  End users can install from{" "}
                  <b>Fleet Desktop {">"} Self-service</b>.
                </>
              }
            >
              Self-service
            </TooltipWrapper>
          </Checkbox>
          <AddPackageAdvancedOptions
            errors={{
              preInstallQuery: formValidation.preInstallQuery?.message,
              postInstallScript: formValidation.postInstallScript?.message,
            }}
            preInstallQuery={formData.preInstallQuery}
            postInstallScript={formData.postInstallScript}
            onChangePreInstallQuery={onChangePreInstallQuery}
            onChangeInstallScript={onChangeInstallScript}
            onChangePostInstallScript={onChangePostInstallScript}
            installScript={formData.installScript}
          />
          <div className="modal-cta-wrap">
            <Button type="submit" variant="brand" disabled={isSubmitDisabled}>
              {isEditingSoftware ? "Save" : "Add software"}
            </Button>
            <Button onClick={onCancel} variant="inverse">
              Cancel
            </Button>
          </div>
        </form>
      )}
    </div>
  );
};

export default AddPackageForm;
