import React from "react";
import classnames from "classnames";

import getMatchedSoftwareIcon from "../";

const baseClass = "software-icon";

type SoftwareIconSizes = "small" | "medium" | "large";

interface ISoftwareIconProps {
  name?: string;
  source?: string;
  size?: SoftwareIconSizes;
  /** Accepts an image url to display for a the software icon image. */
  url?: string;
}

const SOFTWARE_ICON_SIZES: Record<SoftwareIconSizes, string> = {
  small: "24",
  medium: "64",
  large: "96",
};

const SoftwareIcon = ({
  name = "",
  source = "",
  size = "small",
  url,
}: ISoftwareIconProps) => {
  // If we are given a url to render as the icon, we need to render it
  // differently than the svg icons. We will use an img tag instead with the
  // src set to the url.
  if (url) {
    const imgClasses = classnames(
      `${baseClass}__software-img`,
      `${baseClass}__software-img-${size}`
    );
    return (
      <div className={baseClass}>
        <img className={imgClasses} src={url} alt="" />
      </div>
    );
  }

  const MatchedIcon = getMatchedSoftwareIcon({ name, source });
  return (
    <MatchedIcon
      width={SOFTWARE_ICON_SIZES[size]}
      height={SOFTWARE_ICON_SIZES[size]}
      viewBox="0 0 32 32"
      className={baseClass}
    />
  );
};

export default SoftwareIcon;
