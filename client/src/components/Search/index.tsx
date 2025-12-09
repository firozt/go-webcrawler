import './index.css'

type Props = {
	val: string
	setVal: (newVal: string) => void
	errored?: boolean
	handleSubmit?: () => void
	inputTitle?: string
	buttonText?: string
	buttonClickable?: boolean
	placeholder?: string
}

const index = ({val,setVal, errored=false, handleSubmit,inputTitle="",buttonText="",buttonClickable=true,placeholder=""}: Props) => {
  return (
		<div 
			style={{
				outline:`${errored ? "1px solid rgb(233, 92, 92)": "1px solid gray"}`,
			}}
			className='search'>
			{
				inputTitle.length > 0?
				<h2>{inputTitle}</h2> :
				<></>
			}
			<div id='divider'></div>
			<input value={val} onChange={(e) => setVal(e.target.value)} type='search' placeholder={placeholder}/>
			{
				buttonText.length > 0 &&
				<button style={!buttonClickable ? {cursor:"not-allowed", backgroundColor:"black",color:"white"} : {cursor:"pointer"}} onClick={() => buttonClickable && handleSubmit ? handleSubmit(): () => 1}>{buttonText}</button>
			}
		</div>
  )
}

export default index