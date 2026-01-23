import { useState } from 'react'
import './index.css'
import axios, { type AxiosResponse } from 'axios'
import Dropdown from '../../components/DropDown'
import Error from '../../components/Error'
import Search from '../../components/Search'
import RelationshipGraph from '../../components/RelationshipGraph'

type Props = {
	urlInput: string
  toggleShowSearch: () => void
}

type Page = {
  url: string
  title: string
  content: string
}


const SearchPage = ({urlInput,toggleShowSearch}: Props) => {
	const [searchInput, setSearchInput] = useState<string>('')
	const [error, setError] = useState<string>('')
	const [searchResults, setSearchResults] = useState<Page[][]>([])
	const [lastPhrase, setLastPhrase] = useState<string>('')

	
	const handleSearch = () => {
    setError("")
    if (searchInput.length == 0){
      setError("Input length cannot be zero. Please put atleast a single character")
      return
    }
    setLastPhrase(searchInput)
    const API_URL: string = `
    ${import.meta.env.VITE_API_DOMAIN}/${import.meta.env.VITE_API_VER}/search?q=${searchInput}&limit=10&domain=${urlInput}&domain=${urlInput}`
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
		<>
    <RelationshipGraph url={urlInput}/>
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
            <button id='back' onClick={toggleShowSearch}>Back</button>
          {
            error.length > 0 && <Error message={error}/>
          }
          </div>
          <div style={{width:"min(95%,1000px)", margin:"auto",display:"flex",flexDirection:"column",gap:"10px", marginBottom:"2rem"}}>
            {
              searchResults.map((pageGroup, idx) => {
                return (
                  <div key={idx}>
                    <Dropdown title={pageGroup[0].title} content={pageGroup} highlightWord={lastPhrase}/>
                  </div>
                )
              })
            }
          </div>
        </div>  
		</>
  )
}

export default SearchPage