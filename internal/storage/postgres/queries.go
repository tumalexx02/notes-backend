package postgres

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
)

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
	getNoteNodesQuery = `
		SELECT * FROM note_nodes
		WHERE note_id = $1
		ORDER BY "order";
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
)
