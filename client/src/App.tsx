import { useEffect, useState } from 'react'
import './App.css'
import NavBar from './components/NavBar'
import Search from './components/Search'
import Error from './components/Error'
import axios, { type AxiosResponse } from 'axios'
import Spinner from './components/Spinner'

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

function App() {
  const [urlInput, setUrlInput] = useState<string>('')
  const [searchInput, setSearchInput] = useState<string>('')
  const [error, setError] = useState<string>('')
  const [isLightMode, setIsLightMode] = useState<boolean>(false) // light mode state
  const [searchMode, setSearchMode] = useState<boolean>(false) // determines what page to show
  const [searchResults, setSearchResults] = useState<Page[]>([])
  const [buttonClickable, setButtonClickable] = useState<boolean>(true)

  useEffect(() => {
    // sets storage light mode
    setIsLightMode(checkDarkModeStorage())
  },[])

  // for updating lightmode
  useEffect(() => {
    const bg = isLightMode ? "white" : "#242424";
    document.documentElement.style.setProperty("--background-color", bg);
  }, [isLightMode]);

  const handleSearch = () => {
    const API_URL: string = `
    ${import.meta.env.VITE_API_DOMAIN}:${import.meta.env.VITE_API_PORT}/api/${import.meta.env.VITE_API_VER}/search?q=${searchInput}&limit=10`
    console.log(API_URL)
    axios.get(API_URL)
    .then((resp: AxiosResponse<Page[]>) => {
      // obtain data and parse it
      const pageList = resp.data ?? []

      const parsedPageList: Page[] = []
      pageList.forEach((page) => {
        console.warn("page: " + page)
        const parsedText: string[] = getClosestWords(searchInput,page.content,20)
        parsedText.forEach((content) => {
            parsedPageList.push({
              url: page.url,
              title: page.title,
              content: content,
            })
        })
      })
      
      setSearchResults(parsedPageList)
    })
    .catch(err => {
      alert("ERROR : " + err) 
    })

  }

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
    setButtonClickable(false)
    if (!isValidUrl(urlInput)) {
      setError("The URL provided is not a valid http url. Please enter a valid URL to scrap in the form http://www.domain.com")
      setButtonClickable(true)
    }

    const API_URL: string = `${import.meta.env.VITE_API_DOMAIN}:${import.meta.env.VITE_API_PORT}/api/${import.meta.env.VITE_API_VER}/crawl`
    console.warn(API_URL)
    const requestBody: CrawlPostBody = {
      url: urlInput,
      maxDepth: 5,
      followExternal: false
    }
    axios.post(API_URL,requestBody)
    .then(() => {
      setSearchResults([])
      setSearchMode(true) // load search page
      
    }).catch(((err: unknown) => {
      setError("Server Error")
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

const getClosestWords = (keyword: string, largeText: string, windowSize = 20): string[] => {
    const res: string[] = []
    let largeTextArr = largeText.split(/\s+/)

    let indexOfKeyword = largeTextArr.findIndex(word => word === keyword)

    while (indexOfKeyword !== -1) {
        const halfWindow = Math.floor(windowSize / 2)

        const lIndex = Math.max(0, indexOfKeyword - halfWindow)
        const rIndex = Math.min(largeTextArr.length, indexOfKeyword + halfWindow + 1)
        res.push(largeTextArr.slice(lIndex, rIndex).join(" "))

        largeTextArr = [
            ...largeTextArr.slice(0, indexOfKeyword),
            ...largeTextArr.slice(indexOfKeyword + 1)
        ]

        indexOfKeyword = largeTextArr.findIndex(word => word === keyword)
    }

    return res
}

  return (
    <div className={isLightMode ? 'lightmode' : ''}>
      <NavBar isLightMode={isLightMode} toggleLightMode={toggleLightMode}/>
      
      {
        searchMode ?
        <div className='search-page'>
          <div style={{width:"fit-content",margin:"auto",marginBottom:"2rem",display:"flex",flexDirection:"row",gap:"10px"}}>
            <Search
            val={searchInput}
            handleSubmit={handleSearch}
            setVal={(newVal: string) => setSearchInput(newVal)}
            buttonText='Search'
            placeholder='keywords'
            />
            <button onClick={() => setSearchMode(false)}>Back</button>

          </div>
          {
            searchResults.map((page, idx) => {
              return (
                <div className='page-result' key={`${idx}-${page.url}`}>
                  <h3>
                    <a href={page.url} target='_BLANK'>
                      {page.title}
                    </a>
                  </h3>
                  <p>
                    {page.content}
                  </p>
                </div>
              )
            })
          }
        </div>        
        :
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
          handleSubmit={handleSubmit}
          buttonText='Crawl'
          placeholder='https://www.example.com'
          />
          <Error message={error} />
      </div>
      
      }
      {
        !buttonClickable ? <Spinner isLightMode={isLightMode}/> : <div style={{height:"36px"}}></div>
      }
      
    </div>
  )
}

export default App
