-- name: GetSystem :one
SELECT id, "systemName", "localGovernmentId", "createdAt", "updatedAt", 
       "mailAddress", telephone, remark
FROM public.system
WHERE id = $1 LIMIT 1;

-- name: GetSystems :many
SELECT id, "systemName", "localGovernmentId", "createdAt", "updatedAt", 
       "mailAddress", telephone, remark
FROM public.system
ORDER BY "createdAt" DESC;

-- name: GetSystemsByLocalGovernment :many
SELECT id, "systemName", "localGovernmentId", "createdAt", "updatedAt", 
       "mailAddress", telephone, remark
FROM public.system
WHERE "localGovernmentId" = $1
ORDER BY "createdAt" DESC;

-- name: GetSystemByName :one
SELECT id, "systemName", "localGovernmentId", "createdAt", "updatedAt", 
       "mailAddress", telephone, remark
FROM public.system
WHERE "systemName" = $1 LIMIT 1;

-- name: GetSystemsByEmail :many
SELECT id, "systemName", "localGovernmentId", "createdAt", "updatedAt", 
       "mailAddress", telephone, remark
FROM public.system
WHERE "mailAddress" = $1
ORDER BY "createdAt" DESC;

-- name: CreateSystem :one
INSERT INTO public.system ("systemName", "localGovernmentId", "mailAddress", telephone, remark)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, "systemName", "localGovernmentId", "createdAt", "updatedAt", 
          "mailAddress", telephone, remark;

-- name: UpdateSystem :one
UPDATE public.system
SET "systemName" = $2, "localGovernmentId" = $3, "mailAddress" = $4, 
    telephone = $5, remark = $6, "updatedAt" = now()
WHERE id = $1
RETURNING id, "systemName", "localGovernmentId", "createdAt", "updatedAt", 
          "mailAddress", telephone, remark;

-- name: UpdateSystemContact :one
UPDATE public.system
SET "mailAddress" = $2, telephone = $3, "updatedAt" = now()
WHERE id = $1
RETURNING id, "systemName", "localGovernmentId", "createdAt", "updatedAt", 
          "mailAddress", telephone, remark;

-- name: DeleteSystem :exec
DELETE FROM public.system
WHERE id = $1; 