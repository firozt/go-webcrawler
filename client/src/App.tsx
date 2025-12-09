import { useEffect, useState } from 'react'
import './App.css'
import NavBar from './components/NavBar'
import Search from './components/Search'
import Error from './components/Error'

function App() {
  const [urlInput, setUrlInput] = useState<string>('')
  const [error, setError] = useState<string>('')
  const [isLightMode, setIsLightMode] = useState<boolean>(false) // Light mode state


  useEffect(() => {
    // sets storage light mode
    setIsLightMode(checkDarkModeStorage())
  },[])

  // 
  useEffect(() => {
    const bg = isLightMode ? "white" : "#242424";
    document.documentElement.style.setProperty("--background-color", bg);
  }, [isLightMode]);


  // checks if darkmode was set in local storage
  const checkDarkModeStorage = (): boolean => {
    try {
      const storageVal: string = localStorage.getItem("isLightMode") ?? ""
      if (storageVal.toLowerCase() != "true" && storageVal.toLowerCase() != "false") {
        throw Error
      }
      return JSON.parse(storageVal)
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    } catch (_) {
      localStorage.setItem("isLightMode","false")
      return false
    }
  }

  // toggles lightmode bool and saves new value to storage
  const toggleLightMode = () => {
  setIsLightMode(prev => {
    const next = !prev;
    localStorage.setItem("isLightMode", JSON.stringify(next));
    return next;
  });
};

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
    <div className={isLightMode ? 'lightmode' : ''}>
      <NavBar isLightMode={isLightMode} toggleLightMode={toggleLightMode}/>
      <div className='initial-input-page'>
        <h1>Domain Search</h1>
        <p>
          A tool that crawls a domain and indexes all its pages to allows you to quickly search for keywords across the site’s content.
          Documentation for the API and project can be found <a href='https://github.com/firozt/go-webcrawler/blob/main/README.md'><span>here</span></a>
        </p>
        <Search 
        inputTitle='Site'
        val={urlInput} 
        setVal={(newVal: string) => setUrlInput(newVal)} 
        errored={error.length > 0}
        handleSubmit={handleSubmit}
        />
        <Error message={error} />
      </div>
    </div>
  )
}

export default App
