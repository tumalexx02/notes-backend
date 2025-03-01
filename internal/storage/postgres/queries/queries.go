package queries

// note nodes' queries
const (
	CreateBlankNoteNodeQuery = `
		INSERT INTO note_nodes (note_id, "order", content_type, content) 
		VALUES ($1, (SELECT COUNT(*) FROM note_nodes WHERE note_id = $1), $2, $3)
		RETURNING id;
	`
	DeleteNoteNodeQuery = `
		DELETE FROM note_nodes
		WHERE id = $1
		RETURNING note_id, "order";
	`
	UpdateOrderAfterDeleteQuery = `
		UPDATE note_nodes
		SET "order" = "order" - 1
		WHERE note_id = $1 AND "order" > $2;
	`
	UpdateOrderQuery = `
		UPDATE note_nodes
		SET "order" = CASE
				WHEN "order" > $2 AND "order" <= $3 THEN "order" - 1
				WHEN "order" < $2 AND "order" >= $3 THEN "order" + 1
				WHEN "order" = $2 THEN $3
				ELSE "order"
		END
		WHERE note_id = $1;
	`
	UpdateNoteNodeContentQuery = `
		UPDATE note_nodes
		SET content = $2
		WHERE id = $1
		RETURNING note_id;
	`
	NodesCountQuery = `
		SELECT COUNT(*) 
		FROM note_nodes 
		WHERE note_id = $1;
	`
	GetNoteNodeByOrderQuery = `
		SELECT COUNT(*) 
		FROM note_nodes 
		WHERE note_id = $1 AND "order" = $2;
	`
	GetNoteNodesQuery = `
		SELECT * FROM note_nodes
		WHERE note_id = $1
		ORDER BY "order";
	`
	IsUserNodeOwnerQuery = `
    SELECT COUNT(*) FROM note_nodes
    WHERE id = $2 AND (note_id IN (SELECT id FROM notes WHERE user_id = $1));
	`
	GetNoteIdByNoteNodeIdQuery = `
		SELECT * FROM note_nodes
		WHERE id = $1;
	`
	GetAllNotesNodesQuery = `
		SELECT * FROM note_nodes
		WHERE note_id = $1;
	`
)

// notes' queries
const (
	CreateNoteQuery = `
		INSERT INTO notes (title, user_id) 
		VALUES ($1, $2)
		RETURNING id;
	`
	GetNoteQuery = `
		SELECT * FROM notes
		WHERE id = $1;
	`
	GetPublicNoteQuery = `
		SELECT * FROM notes
		WHERE id = $1 AND is_public = TRUE;
	`
	GetNotesByUserIdQuery = `
		SELECT * FROM notes
		WHERE user_id = $1
		ORDER BY updated_at DESC;
	`
	UpdateNoteTitleQuery = `
		UPDATE notes
		SET title = $2, updated_at = NOW()
		WHERE id = $1;
	`
	ArchiveNoteQuery = `
		UPDATE notes
		SET archived_at = NOW()
		WHERE id = $1;
	`
	UnarchiveNoteQuery = `
		UPDATE notes
		SET archived_at = NULL
		WHERE id = $1;
	`
	DeleteNoteQuery = `
		DELETE FROM notes
		WHERE id = $1;
	`
	SetUpdatedAtQuery = `
		UPDATE notes
		SET updated_at = NOW()
		WHERE id = $1;
	`
	IsUserNoteOwnerQuery = `
		SELECT COUNT(*) FROM notes
		WHERE id = $2 AND user_id = $1;`
	MakeNotePublicQuery = `
		UPDATE notes
		SET is_public = TRUE, updated_at = NOW()
		WHERE id = $1;
	`
	MakeNotePrivateQuery = `
		UPDATE notes
		SET is_public = FALSE, updated_at = NOW()
		WHERE id = $1;
	`
)

// users' queries
const (
	CreateUserQuery = `
		INSERT INTO users (email, name, password_hash) 
		VALUES ($1, $2, $3)
		RETURNING id;
	`
	GetUserByEmailQuery = `
		SELECT * FROM users
		WHERE email = $1;
	`
	GetUserByIdQuery = `
		SELECT * FROM users
		WHERE id = $1;
	`
)

// auth tokens' queries
const (
	CreateRefreshTokenQuery = `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at) 
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	GetRefreshTokenByIdQuery = `
		SELECT * FROM refresh_tokens
		WHERE id = $1;
	`
	DeleteExpiredRefreshTokensQuery = `
		DELETE FROM refresh_tokens
		WHERE expires_at < NOW();
	`
	DeleteRefreshTokenQuery = `
		DELETE FROM refresh_tokens
		WHERE id = $1;
	`
)
