CREATE TABLE public.user (
    id SERIAL PRIMARY KEY,
    name character varying NOT NULL,
    email character varying NOT NULL UNIQUE,
    password character varying NOT NULL,
    hash character varying UNIQUE,
    zeroTierNetworkId character varying,
    zeroTierDiscoIP character varying,
    homeLocation jsob,
    deployStatus integer NOT NULL DEFAULT 0,
    deployIP character varying,
    deployID character varying,
    defaultRegion character varying,
    deployedAt timestamp without time zone
);