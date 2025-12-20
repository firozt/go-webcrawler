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
};

const SearchPage = ({ title, content = [], highlightWord }: Props) => {
  const [showDropDown, setShowDropDown] = useState<boolean>(false);

  return (
    <div className="dropdown">
      <div
        className="header"
        onClick={() => setShowDropDown(prev => !prev)}
      >
        <div className="button">
          <h3>{title}</h3>
          <p id="entrycount">{content.length} entries</p>
        </div>

        <div>
          {showDropDown ? null : (
            <img
              width={15}
              src="/arrowdown.svg"
              className="arrow"
            />
          )}
        </div>
      </div>

      <div className={`dropdown-content ${showDropDown ? 'open' : ''}`}>
        {content.map((page, idx) => {
          const regex = new RegExp(`(${highlightWord})`, 'gi');
          const parts = page.content.split(regex);

          return (
            <div className="page-result" key={`${idx}-${page.url}`}>
              <h3>
                <a href={page.url} target="_BLANK" rel="noreferrer">
                  {page.url}
                </a>
              </h3>
              <p>
                {parts.map((part, i) =>
                  regex.test(part) ? (
                    <span key={i} className="highlight">
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

export default SearchPage;
