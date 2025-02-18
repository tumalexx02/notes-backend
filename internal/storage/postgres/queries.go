package postgres

// note nodes' queries
const (
	createBlankNoteNodeQuery = `
		INSERT INTO note_nodes (note_id, "order", content_type, content) 
		VALUES ($1, (SELECT COUNT(*) FROM note_nodes WHERE note_id = $1), $2, $3)
		RETURNING id;
	`
	deleteNoteNodeQuery = `
		DELETE FROM note_nodes
		WHERE id = $1
		RETURNING note_id, "order";
	`
	updateOrderAfterDeleteQuery = `
		UPDATE note_nodes
		SET "order" = "order" - 1
		WHERE note_id = $1 AND "order" > $2;
	`
	updateOrderQuery = `
		UPDATE note_nodes
		SET "order" = CASE
				WHEN "order" > $2 AND "order" <= $3 THEN "order" - 1
				WHEN "order" < $2 AND "order" >= $3 THEN "order" + 1
				WHEN "order" = $2 THEN $3
				ELSE "order"
		END
		WHERE note_id = $1;
	`
	updateNoteNodeContentQuery = `
		UPDATE note_nodes
		SET content = $2
		WHERE id = $1
		RETURNING note_id;
	`
	nodesCountQuery = `
		SELECT COUNT(*) 
		FROM note_nodes 
		WHERE note_id = $1;
	`
	getNoteNodeByOrderQuery = `
		SELECT COUNT(*) 
		FROM note_nodes 
		WHERE note_id = $1 AND "order" = $2;
	`
	getNoteNodesQuery = `
		SELECT * FROM note_nodes
		WHERE note_id = $1
		ORDER BY "order";
	`
	isUserNodeOwnerQuery = `
    SELECT COUNT(*) FROM note_nodes
    WHERE id = $2 AND (note_id IN (SELECT id FROM notes WHERE user_id = $1));
	`
)

// notes' queries
const (
	createNoteQuery = `
		INSERT INTO notes (title, user_id) 
		VALUES ($1, $2)
		RETURNING id;
	`
	getNoteQuery = `
		SELECT * FROM notes
		WHERE id = $1;
	`
	getNotesByUserIdQuery = `
		SELECT * FROM notes
		WHERE user_id = $1
		ORDER BY updated_at DESC;
	`
	updateNoteTitleQuery = `
		UPDATE notes
		SET title = $2, updated_at = NOW()
		WHERE id = $1;
	`
	archiveNoteQuery = `
		UPDATE notes
		SET archived_at = NOW()
		WHERE id = $1;
	`
	unarchiveNoteQuery = `
		UPDATE notes
		SET archived_at = NULL
		WHERE id = $1;
	`
	deleteNoteQuery = `
		DELETE FROM notes
		WHERE id = $1;
	`
	setUpdatedAtQuery = `
		UPDATE notes
		SET updated_at = NOW()
		WHERE id = $1;
	`
	isUserNoteOwnerQuery = `
		SELECT COUNT(*) FROM notes
		WHERE id = $2 AND user_id = $1;`
)

// users' queries
const (
	createUserQuery = `
		INSERT INTO users (email, name, password_hash) 
		VALUES ($1, $2, $3)
		RETURNING id;
	`
	getUserByEmailQuery = `
		SELECT * FROM users
		WHERE email = $1;
	`
	getUserByIdQuery = `
		SELECT * FROM users
		WHERE id = $1;
	`
)

// auth tokens' queries
const (
	createRefreshTokenQuery = `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at) 
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	getRefreshTokenByIdQuery = `
		SELECT * FROM refresh_tokens
		WHERE id = $1;
	`
	deleteExpiredRefreshTokensQuery = `
		DELETE FROM refresh_tokens
		WHERE expires_at < NOW();
	`
	deleteRefreshTokenQuery = `
		DELETE FROM refresh_tokens
		WHERE id = $1;
	`
)
