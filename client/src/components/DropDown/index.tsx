import { useState } from 'react';
import './index.css';

type Page = {
  url: string;
  title: string;
  content: string;
};

type Props = {
  title: string;
  content: Page[];
  highlightWord: string;
  isLightMode: boolean;
};

const Index = ({ title, content = [], highlightWord, isLightMode }: Props) => {
  const [showDropDown, setShowDropDown] = useState<boolean>(false);
  console.log(isLightMode)
  return (
    <div className='dropdown' style={{borderColor:`${isLightMode ? "black" : "#a1a1a1ff"}`}}>

      <div className='header' onClick={() => setShowDropDown(prev => !prev)} style={{backgroundColor:`${isLightMode ? "#e0e0e0" : "#656464ff"}`}}>
        <div className='button'>
          <h3>{title}</h3>
          <p id='entrycount'>{content.length} entries</p>
        </div>
        <div>{showDropDown ? <p></p> : <img width={15} src={"/arrowdown.svg"} style={ !isLightMode ? {filter:"invert(100)"} : {}} />}</div>
      </div>
      <div className={`dropdown-content ${showDropDown ? 'open' : ''}`}>
        {content.map((page, idx) => {
          const regex = new RegExp(`(${highlightWord})`, 'gi');
          const parts = page.content.split(regex);
          return (
            <div className='page-result' key={`${idx}-${page.url}`}>
              <h3>
                <a href={page.url} target='_BLANK'>
                  {page.url}
                </a>
              </h3>
              <p>
                {parts.map((part, i) =>
                  regex.test(part) ? (
                    <span
                      key={i}
                      style={{
                        backgroundColor: '#e2cdabff',
                        color: '#242424',
                        padding: '1px',
                      }}
                    >
                      {part}
                    </span>
                  ) : (
                    part
                  )
                )}
              </p>
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default Index;
