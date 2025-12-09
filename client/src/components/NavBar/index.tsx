import './index.css'

type Props = {
	isLightMode: boolean
	toggleLightMode: () => void
}

const NavItemContainerStyle:React.CSSProperties = {display:"flex", flexDirection:"row",justifyContent:"center",alignItems:"center", gap:"10px", cursor:"pointer"} 

const index = ({isLightMode,toggleLightMode}: Props) => {
  return (
      <div className='navbar'>
        <div style={NavItemContainerStyle}>
          <a 
          href='https://www.github.com/firozt/go-webcrawler'  
          target='_BLANK'>
            <img
            width={"40px"}
            id='github'
            src={isLightMode ? '/lightmode-github.png' : '/darkmode-github.png'}
            />
          </a>
          <p>firozt/go-webcrawler</p>
        </div>
        <div>
          <img 
          onClick={toggleLightMode} 
          style={{filter:`${!isLightMode? "invert(100)":"" }`}} 
          id='mode' 
          width={40}
          height={40}
          src={isLightMode ? "/darkmode.png" : "/lightmode.svg"}
          />
        </div>

      </div>
  )
}

export default index