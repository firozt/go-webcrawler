CREATE VIRTUAL TABLE pages USING fts5(
    url,
    title,
    content,
    crawled_at,
    tokenize='porter'
);

