CREATE TABLE public.user (
    id SERIAL PRIMARY KEY,
    name character varying NOT NULL,
    email character varying NOT NULL UNIQUE,
    password character varying NOT NULL,
    hash character varying UNIQUE,
    zeroTierNetworkId character varying,
    zeroTierDiscoIP character varying,
    homeLocation jsonb,
    deployStatus integer NOT NULL DEFAULT 0,
    deployIP character varying,
    deployPassword character varying,
    deployID character varying,
    deployRegion character varying,
    deployedAt timestamp without time zone,
    share_location boolean DEFAULT FALSE
);