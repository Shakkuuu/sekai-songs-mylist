CREATE TABLE my_lists (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    position INT NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE my_list_charts (
    id SERIAL PRIMARY KEY,
    my_list_id INT REFERENCES my_lists(id),
    chart_id INT REFERENCES charts(id),
    clear_type INT,
    memo TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE my_list_chart_attachments (
    id SERIAL PRIMARY KEY,
    my_list_chart_id INT REFERENCES my_list_charts(id),
    attachment_type INT,
    file_url TEXT,
    caption TEXT,
    created_at TIMESTAMP
);
