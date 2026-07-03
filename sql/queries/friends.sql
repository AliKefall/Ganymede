-- name: CreateFriendship :exec
INSERT INTO friendships(user_id, friend_id, created_at)
VALUES($1, $2, $3);

-- name: FriendshipExists :one
SELECT COUNT(1)
FROM friendships
WHERE user_id = $1 AND friend_id = $2;

-- name: ListFriendsByUserID :many
SELECT u.*
FROM friendships AS f
JOIN users AS u ON u.id = CASE
WHEN f.user_id = $1 THEN f.friend_id
ELSE f.user_id
END
WHERE f.user_id = $2 OR f.friend_id = $3
ORDER BY u.username;

-- name: DeleteFriendship :exec
DELETE FROM friendships
WHERE user_id = $1 AND friend_id = $2;

-- name: CreateFriendRequest :exec
INSERT INTO friend_requests(requester_id, target_id, created_at)
VALUES($1, $2, $3);

-- name: FriendRequestExists :one
SELECT COUNT(1)
FROM friend_requests
WHERE (requester_id = $1 AND target_id = $2) OR (requester_id = $2 AND target_id = $1);

-- name: DeleteFriendRequest :exec
DELETE FROM friend_requests
WHERE (requester_id = $1 AND target_id = $2) OR (requester_id = $2 AND target_id = $1);

-- name: ListIncomingFriendRequestsByUserID :many
SELECT u.*
FROM friend_requests AS fr
JOIN users AS u ON u.id = fr.requester_id
WHERE fr.target_id = $1
ORDER BY fr.created_at DESC;

-- name: ListOutgoingFriendRequestsByUserID :many
SELECT u.*
FROM friend_requests AS fr
JOIN users AS u ON u.id = fr.target_id
WHERE fr.requester_id = $1
ORDER BY fr.created_at DESC;
