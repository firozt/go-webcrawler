import Search from "../../components/Search"
import Error from '../../components/Error'
import { useState } from "react"
import axios from "axios"
import Spinner from "../../components/Spinner"
import './index.css'

type Props = {
	urlInput: string
	setUrlInput: (newVal: string) => void
	setShowSearch: () => void
}

type CrawlPostBody = {
  url: string
  maxDepth: number
  followExternal: boolean
}

type Page = {
  url: string
  title: string
  content: string
}


const Page = ({urlInput, setUrlInput, setShowSearch}: Props) => {

	const [error, setError] = useState<string>('') 
	const [buttonClickable, setButtonClickable] = useState<boolean>(true)
	
	const handleSubmit = () => {
    setError("")
    setButtonClickable(false)
    if (!isValidUrl(urlInput)) {
      setError("The URL provided is not a valid http url. Please enter a valid URL to scrap in the form http://www.domain.com")
      setButtonClickable(true)
      return
    }

    const API_URL: string = `${import.meta.env.VITE_API_DOMAIN}/api/${import.meta.env.VITE_API_VER}/crawl`
    console.warn(API_URL)
    const requestBody: CrawlPostBody = {
      url: urlInput,
      maxDepth: 5,
      followExternal: false
    }
    axios.post(API_URL,requestBody)
    .then(() => {
			setShowSearch()
    }).catch(((err: unknown) => {
      setError("Server Error, please try again shortly")
      console.error(err)
    })).finally(() => setButtonClickable(true))
    
  }

  const isValidUrl = (input: string): boolean => {
    try {
      new URL(input); // tries to parse the string as a URL
      return true;
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    } catch (_) {
      return false;
    }
  }

  return (
		<>
			<div className='initial-input-page'>
				<h1>Domain Search</h1>
					<p>
            A tool that crawls a domain and indexes all its pages to allows you to quickly search for keywords across the site’s content.
            Documentation for the API and project can be found <a href='https://github.com/firozt/go-webcrawler/blob/main/README.md'><span>here</span></a>
          </p>
          <Search 
          buttonClickable={buttonClickable}
          inputTitle='Site'
          val={urlInput} 
          setVal={(newVal: string) => setUrlInput(newVal)} 
          errored={error.length > 0}
          handleSubmit={() => handleSubmit()}
          buttonText='Crawl'
          placeholder='https://www.example.com'
          />
          <Error message={error} />
      </div>
			{
        !buttonClickable ? <Spinner subtext={`Crawling ${urlInput.split("/")[2]}`}/> : <div style={{height:"36px"}}></div>
			}
		</>
  )
}

export default Page