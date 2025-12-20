import "./index.css";

// const LightModeHex = "#242424"
// const DarkModeHex = "#a25757ff"

type Props = {
  subtext?: string
}

const index = ({subtext=""}: Props) => {
  return (
    <div className="spinner-container">
      <div className="spinner"></div>
      <p>{subtext}</p>
    </div>
  );
};

export default index;
