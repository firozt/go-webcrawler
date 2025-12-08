import { useState } from 'react'
import './App.css'

function App() {
  const [urlInput, setUrlInput] = useState<string>('')
  
  const [error, setError] = useState<string>('')

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setUrlInput(e.target.value)
  }

  const handleSubmit = () => {
    setError("")
    if (!isValidUrl(urlInput)) {
      setError("The URL provided is not a valid http url. Please enter a valid URL to scrap in the form http://www.domain.com")
    }
    // to http req
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
      <div className='navbar'>
        <p>firozt/DOMAIN SEARCH</p>
        <a href='https://www.github.com/firozt/go-webcrawler' target='_BLANK'>
          <img id='github' src='../public/image.png'/>
        </a>
      </div>
      <div className='page'>
        <h1>Domain Search</h1>
        <p>
          A tool crawls a website and indexes all its pages. It allows you to quickly search for keywords across the site’s content.
          Documentation for the API and projet can be found <a href=''><span>here</span></a>
        </p>
        <div 
        style={{
          outline:`${error.length > 0 &&"1px solid rgb(233, 92, 92)"}`
        }}
        className='search'>
          <h2>Site</h2>
          <div id='divider'></div>
          <input onChange={(e) => handleChange(e)} type='search' placeholder='https://www.example.com'/>
          <button onClick={handleSubmit}>Crawl</button>
        </div>
        {
          error.length > 0 && 
          <>
            <div className='error'>
              <p>{error}</p>
            </div>
          </>
        }
      </div>
    </>
  )
}

export default App
