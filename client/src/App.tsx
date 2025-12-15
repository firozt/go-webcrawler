import { useEffect, useState } from 'react'
import './App.css'
import NavBar from './components/NavBar'
import Search from './components/Search'
import Error from './components/Error'
import axios, { type AxiosResponse } from 'axios'
import Spinner from './components/Spinner'
import Dropdown from './components/DropDown'

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
  const [searchResults, setSearchResults] = useState<Page[][]>([])
  const [buttonClickable, setButtonClickable] = useState<boolean>(true)
  const [lastPhrase, setLastPhrase] = useState<string>('')

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
    setError("")
    if (searchInput.length == 0){
      setError("Input length cannot be zero. Please put atleast a single character")
      return
    }
    setLastPhrase(searchInput)
    const API_URL: string = `
    ${import.meta.env.VITE_API_DOMAIN}/api/${import.meta.env.VITE_API_VER}/search?q=${searchInput}&limit=10&domain=${urlInput}&domain=${urlInput}`
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
      
      setSearchResults(groupPagesByUrl(parsedPageList))
    })
    .catch(err => {
      setError("Server Error, please try again later.:"+ err)
    })
  }


  const groupPagesByUrl = (pages: Page[]): Page[][] => {
  const map = new Map<string, Page[]>();

  for (const page of pages) {
    if (!map.has(page.url)) {
      map.set(page.url, []);
    }
    map.get(page.url)!.push(page);
  }

  return Array.from(map.values());
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
    const words = largeText.split(/\s+/)
    const halfWindow = Math.floor(windowSize / 2)
    let startIndex = 0

    while (true) {
        // find the next occurrence of the keyword as a substring
        const idx = largeText.indexOf(keyword, startIndex)
        if (idx === -1) break

        // count words before the match
        const preWords = largeText.slice(0, idx).split(/\s+/)

        const lIndex = Math.max(0, preWords.length - halfWindow)
        const rIndex = Math.min(words.length, preWords.length + keyword.split(/\s+/).length + halfWindow)

        res.push(words.slice(lIndex, rIndex).join(" "))

        // move startIndex past this match
        startIndex = idx + keyword.length
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
            errored={error.length > 0}
            />
            <button onClick={() => setSearchMode(false)}>Back</button>
          {
            error.length > 0 && <Error message={error}/>
          }
          </div>
          <div style={{width:"min(95%,1000px)", margin:"auto",display:"flex",flexDirection:"column",gap:"10px"}}>
            {
              searchResults.map((pageGroup, idx) => {
                return (
                  <div key={idx}>
                    <Dropdown isLightMode={isLightMode} title={pageGroup[0].title} content={pageGroup} highlightWord={lastPhrase}/>
                  </div>
                )
              })
            }
          </div>
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
