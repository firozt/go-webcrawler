import './index.css'

type Props = {
	val: string
	setVal: (newVal: string) => void
	errored?: boolean
	handleSubmit?: () => void
	inputTitle?: string
	
}

const index = ({val,setVal, errored=false, handleSubmit,inputTitle=""}: Props) => {
  return (
		<div 
			style={{
				outline:`${errored ? "1px solid rgb(233, 92, 92)": "1px solid gray"}`,
				border: "1px solid black"
			}}
			className='search'>
			{
				inputTitle.length > 0?
				<h2>Site</h2> :
				<></>
			}
			<div id='divider'></div>
			<input value={val} onChange={(e) => setVal(e.target.value)} type='search' placeholder='https://www.example.com'/>
			<button onClick={handleSubmit}>Crawl</button>
		</div>
  )
}

export default index