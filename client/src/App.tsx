import { useEffect, useState } from 'react'
import './App.css'
import NavBar from './components/NavBar'
import Error from './components/Error'
import Home from './pages/Home'
import SearchPage from './pages/Search'

function App() {
  const [urlInput, setUrlInput] = useState<string>('')
  const [isLightMode, setIsLightMode] = useState<boolean>(false) // light mode state
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
  
  return (
    <div className={isLightMode ? 'lightmode' : ''}>
      <NavBar 
      isLightMode={isLightMode} 
      toggleLightMode={toggleLightMode}
      />
      { 
        showSearch ?
        <SearchPage
        urlInput={urlInput}
        toggleShowSearch={() => setShowSearch(prev => !prev)}
        />
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
