UPDATE wp_posts SET post_content = REPLACE (post_content, 'src="https://news.china.com.au//wp-content/', 'src="https://cdn.china.com.au/wp-content/') where id = 69713;
UPDATE wp_posts SET post_content = REPLACE (post_content, 'src="https://news.china.com.au/wp-content/', 'src="https://cdn.china.com.au/wp-content/') where id = 69713;


-- replace post_content of wp_posts
UPDATE wp_posts SET post_content = REPLACE (post_content, 'src="https://news.china.com.au//wp-content/', 'src="https://cdn.china.com.au/wp-content/');
UPDATE wp_posts SET post_content = REPLACE (post_content, 'src="https://news.china.com.au/wp-content/', 'src="https://cdn.china.com.au/wp-content/');

-- replace guid of wp_posts
UPDATE wp_posts SET guid = REPLACE (guid, 'https://news.china.com.au//wp-content', 'https://cdn.china.com.au/wp-content') WHERE post_type = 'attachment';
UPDATE wp_posts SET guid = REPLACE (guid, 'https://news.china.com.au/wp-content', 'https://cdn.china.com.au/wp-content') WHERE post_type = 'attachment';

-- replace post_content of wp_posts
UPDATE wp_posts SET post_content = REPLACE (post_content, 'https://news.china.com.au//wp-content', 'https://cdn.china.com.au/wp-content') WHERE post_type = 'attachment';
UPDATE wp_posts SET post_content = REPLACE (post_content, 'https://news.china.com.au/wp-content', 'https://cdn.china.com.au/wp-content') WHERE post_type = 'attachment';

-- replace post_content_filtered of wp_posts
UPDATE wp_posts SET post_content_filtered = REPLACE (post_content_filtered, 'src="https://news.china.com.au//wp-content/', 'src="https://cdn.china.com.au/wp-content/');
UPDATE wp_posts SET post_content_filtered = REPLACE (post_content_filtered, 'src="https://news.china.com.au/wp-content/', 'src="https://cdn.china.com.au/wp-content/');

-- https://myqqjd.com/1579.html
-- https://www.ypojie.com/5847.html
-- https://www.wpzzq.com/1019.html
-- https://cloud.tencent.com/developer/article/1553652