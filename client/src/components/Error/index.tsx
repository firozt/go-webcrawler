import './index.css'

type Props = {
	message: string
	showMessage?: boolean
}

const index = ({message, showMessage=message.length>0}: Props) => {
	if (!showMessage) {
		return (
			<>
			</>
		)
	} else {
		return (
			<div className='error'>
				<p>{message}</p>
			</div>
		)
	}
}

export default index