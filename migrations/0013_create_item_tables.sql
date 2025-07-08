CREATE TABLE IF NOT EXISTS items (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INTEGER NOT NULL,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    uploaded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    salt VARCHAR(255) NOT NULL,
    iv VARCHAR(255) NOT NULL,
    encryption_metadata TEXT,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    
    INDEX idx_user_id (user_id),
    INDEX idx_type (type),
    INDEX idx_uploaded_at (uploaded_at)
);

CREATE TABLE IF NOT EXISTS text_items (
    item_id INT PRIMARY KEY,
    content TEXT NOT NULL,
    
    FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS file_items (
    item_id INT PRIMARY KEY,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    file_created_at DATETIME,
    file_updated_at DATETIME,
    file_url VARCHAR(500) NOT NULL,
    
    FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
    
    INDEX idx_file_name (file_name),
    INDEX idx_mime_type (mime_type)
);
