CREATE TABLE public.user (
    id SERIAL PRIMARY KEY,
    name character varying NOT NULL,
    email character varying NOT NULL UNIQUE,
    password character varying NOT NULL,
    hash character varying UNIQUE,
    zeroTierNetworkId character varying,
    zeroTierDiscoIP character varying,
    homeLocation jsob,
    deployStatus number NOT NULL DEFAULT 0,
    deployIP character varying,
    deployID character varying,
    deployedAt time without timezone
);