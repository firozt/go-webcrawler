import { useEffect, useState } from 'react'
import './App.css'
import NavBar from './components/NavBar'
import Search from './components/Search'
import Error from './components/Error'
import axios, { type AxiosResponse } from 'axios'
import Dropdown from './components/DropDown'
import Home from './pages/Home'


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
  const [searchResults, setSearchResults] = useState<Page[][]>([])
  const [lastPhrase, setLastPhrase] = useState<string>('')
  const [showSearch, setShowSearch] = useState<boolean>(false)

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
        showSearch ?
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
            <button onClick={() => setShowSearch(prev => !prev)}>Back</button>
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
        <Home
        urlInput={urlInput}
        setUrlInput={setUrlInput}
        setShowSearch={() => setShowSearch(true)}
        />
      }
      
    </div>
  )
}

export default App
