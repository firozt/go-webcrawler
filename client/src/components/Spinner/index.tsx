import "./index.css";

const LightModeHex = "#242424"
const DarkModeHex = "#a25757ff"

type Props = {
    isLightMode: boolean
}

const index = ({isLightMode}:Props) => {
  return (
    <div className="spinner-container">
      <div className="spinner" style={{borderTop:`4px solid ${isLightMode ? LightModeHex : DarkModeHex}`}}></div>
    </div>
  );
};

export default index;
