-- OIDC login replaces password authentication. Existing rows keep their
-- hashes until their owners sign in through the identity provider, at which
-- point the identity columns are linked by verified email.
alter table "user" alter column password_hash drop not null;

alter table "user" add column oidc_issuer text;
alter table "user" add column oidc_subject text;

create unique index user_oidc_identity_idx
    on "user" (oidc_issuer, oidc_subject)
    where oidc_issuer is not null;
