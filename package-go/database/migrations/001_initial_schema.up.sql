--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: drizzle; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA drizzle;

--
-- Name: __drizzle_migrations; Type: TABLE; Schema: drizzle; Owner: -
--

CREATE TABLE drizzle.__drizzle_migrations (
    id integer NOT NULL,
    hash text NOT NULL,
    created_at bigint
);

--
-- Name: __drizzle_migrations_id_seq; Type: SEQUENCE; Schema: drizzle; Owner: -
--

ALTER TABLE drizzle.__drizzle_migrations ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME drizzle.__drizzle_migrations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

--
-- Name: m_localGovernment; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."m_localGovernment" (
    id character varying(6) NOT NULL,
    "prefectureName" character varying(255) NOT NULL,
    "cityName" character varying(255) NOT NULL,
    "prefectureNameKana" character varying(255) NOT NULL,
    "cityNameKana" character varying(255) NOT NULL,
    "createdAt" timestamp with time zone DEFAULT now() NOT NULL,
    "updatedAt" timestamp with time zone DEFAULT now() NOT NULL
);

--
-- Name: m_organizationCategory; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."m_organizationCategory" (
    id integer NOT NULL,
    "organizationCategoryNameJa" character varying(255) NOT NULL,
    "organizationCategoryNameEn" character varying(255) NOT NULL,
    "createdAt" timestamp with time zone DEFAULT now() NOT NULL,
    "updatedAt" timestamp with time zone DEFAULT now() NOT NULL
);

--
-- Name: m_organizationCategory_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public."m_organizationCategory" ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public."m_organizationCategory_id_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

--
-- Name: m_userRole; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."m_userRole" (
    id integer NOT NULL,
    "roleNameJa" character varying(255) NOT NULL,
    "roleNameEn" character varying(255) NOT NULL,
    "createdAt" timestamp with time zone DEFAULT now() NOT NULL,
    "updatedAt" timestamp with time zone DEFAULT now() NOT NULL
);

--
-- Name: m_userRole_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public."m_userRole" ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public."m_userRole_id_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

--
-- Name: gcasUser; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."gcasUser" (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    "familyName" character varying(60) NOT NULL,
    "givenName" character varying(60) NOT NULL,
    "mailAddress" character varying(255) NOT NULL,
    "organizationCategoryId" integer,
    "createdAt" timestamp with time zone DEFAULT now() NOT NULL,
    "updatedAt" timestamp with time zone DEFAULT now() NOT NULL,
    "lastLoginAt" timestamp with time zone
);

--
-- Name: gcasGroup; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."gcasGroup" (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    "groupCategoryId" integer,
    "groupName" character varying(255) NOT NULL,
    "createdAt" timestamp with time zone DEFAULT now() NOT NULL,
    "updatedAt" timestamp with time zone DEFAULT now() NOT NULL
);

--
-- Name: gcasGroupUserRelation; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."gcasGroupUserRelation" (
    "gcasUserId" uuid NOT NULL,
    "groupId" uuid NOT NULL,
    "createdAt" timestamp with time zone DEFAULT now() NOT NULL,
    "updatedAt" timestamp with time zone DEFAULT now() NOT NULL,
    "userRoleId" integer NOT NULL
);

--
-- Name: gcasGroupSystemRelation; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."gcasGroupSystemRelation" (
    "systemId" uuid NOT NULL,
    "groupId" uuid NOT NULL,
    "createdAt" timestamp with time zone DEFAULT now() NOT NULL,
    "updatedAt" timestamp with time zone DEFAULT now() NOT NULL
);

--
-- Name: project; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.project (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    "projectName" character varying(255) NOT NULL,
    "localGovernmentId" character varying(6) NOT NULL,
    "projectType" character varying(255) NOT NULL,
    "governmentCloudConnectionType" character varying(255) NOT NULL,
    "createdAt" timestamp with time zone DEFAULT now() NOT NULL,
    "updatedAt" timestamp with time zone DEFAULT now() NOT NULL,
    "corporateNumber" character varying(13) NOT NULL,
    "vendorName" character varying(255) NOT NULL,
    "serviceOutsourcingFee" integer,
    "cloudUsageFee" integer
);

--
-- Name: projectCost; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."projectCost" (
    "projectId" uuid NOT NULL,
    year integer NOT NULL,
    cost integer,
    "createdAt" timestamp with time zone DEFAULT now() NOT NULL,
    "updatedAt" timestamp with time zone DEFAULT now() NOT NULL
);

--
-- Name: projectSystemRelation; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."projectSystemRelation" (
    "projectId" uuid NOT NULL,
    "systemId" uuid NOT NULL,
    "createdAt" timestamp with time zone DEFAULT now() NOT NULL,
    "updatedAt" timestamp with time zone DEFAULT now() NOT NULL
);

--
-- Name: system; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.system (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    "systemName" character varying(255) NOT NULL,
    "localGovernmentId" character varying(6),
    "createdAt" timestamp with time zone DEFAULT now() NOT NULL,
    "updatedAt" timestamp with time zone DEFAULT now() NOT NULL,
    "mailAddress" character varying(255) NOT NULL,
    telephone character varying(255),
    remark character varying(1000)
);

--
-- Name: systemBasicInformation; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."systemBasicInformation" (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    "projectId" uuid NOT NULL,
    "corporateNumber" character varying(13) NOT NULL,
    "vendorName" character varying(255) NOT NULL,
    "operationStartDate" character varying(255) NOT NULL,
    "standardizationTasks" jsonb NOT NULL,
    "createdAt" timestamp with time zone DEFAULT now() NOT NULL,
    "updatedAt" timestamp with time zone DEFAULT now() NOT NULL
);

--
-- Name: __drizzle_migrations __drizzle_migrations_pkey; Type: CONSTRAINT; Schema: drizzle; Owner: -
--

ALTER TABLE ONLY drizzle.__drizzle_migrations
    ADD CONSTRAINT __drizzle_migrations_pkey PRIMARY KEY (id);

--
-- Name: m_localGovernment m_localGovernment_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."m_localGovernment"
    ADD CONSTRAINT "m_localGovernment_pkey" PRIMARY KEY (id);

--
-- Name: m_organizationCategory m_organizationCategory_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."m_organizationCategory"
    ADD CONSTRAINT "m_organizationCategory_pkey" PRIMARY KEY (id);

--
-- Name: m_userRole m_userRole_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."m_userRole"
    ADD CONSTRAINT "m_userRole_pkey" PRIMARY KEY (id);

--
-- Name: gcasUser gcasUser_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."gcasUser"
    ADD CONSTRAINT "gcasUser_pkey" PRIMARY KEY (id);

--
-- Name: gcasGroup gcasGroup_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."gcasGroup"
    ADD CONSTRAINT "gcasGroup_pkey" PRIMARY KEY (id);

--
-- Name: gcasGroupUserRelation gcasGroupUserRelation_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."gcasGroupUserRelation"
    ADD CONSTRAINT "gcasGroupUserRelation_pkey" PRIMARY KEY ("gcasUserId", "groupId");

--
-- Name: gcasGroupSystemRelation gcasGroupSystemRelation_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."gcasGroupSystemRelation"
    ADD CONSTRAINT "gcasGroupSystemRelation_pkey" PRIMARY KEY ("systemId", "groupId");

--
-- Name: project project_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.project
    ADD CONSTRAINT project_pkey PRIMARY KEY (id);

--
-- Name: projectCost projectCost_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."projectCost"
    ADD CONSTRAINT "projectCost_pkey" PRIMARY KEY ("projectId", year);

--
-- Name: projectSystemRelation projectSystemRelation_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."projectSystemRelation"
    ADD CONSTRAINT "projectSystemRelation_pkey" PRIMARY KEY ("projectId", "systemId");

--
-- Name: system system_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.system
    ADD CONSTRAINT system_pkey PRIMARY KEY (id);

--
-- Name: systemBasicInformation systemBasicInformation_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."systemBasicInformation"
    ADD CONSTRAINT "systemBasicInformation_pkey" PRIMARY KEY (id);

--
-- Name: gcasUser_mailAddress_unique; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX "gcasUser_mailAddress_unique" ON public."gcasUser" USING btree ("mailAddress");

--
-- Name: gcasGroup_groupName_unique; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX "gcasGroup_groupName_unique" ON public."gcasGroup" USING btree ("groupName");

--
-- Name: system_systemName_unique; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX "system_systemName_unique" ON public.system USING btree ("systemName");

--
-- Name: gcasUser gcasUser_organizationCategoryId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."gcasUser"
    ADD CONSTRAINT "gcasUser_organizationCategoryId_fkey" FOREIGN KEY ("organizationCategoryId") REFERENCES public."m_organizationCategory"(id) ON UPDATE CASCADE ON DELETE SET NULL;

--
-- Name: gcasGroupUserRelation gcasGroupUserRelation_gcasUserId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."gcasGroupUserRelation"
    ADD CONSTRAINT "gcasGroupUserRelation_gcasUserId_fkey" FOREIGN KEY ("gcasUserId") REFERENCES public."gcasUser"(id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: gcasGroupUserRelation gcasGroupUserRelation_groupId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."gcasGroupUserRelation"
    ADD CONSTRAINT "gcasGroupUserRelation_groupId_fkey" FOREIGN KEY ("groupId") REFERENCES public."gcasGroup"(id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: gcasGroupUserRelation gcasGroupUserRelation_userRoleId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."gcasGroupUserRelation"
    ADD CONSTRAINT "gcasGroupUserRelation_userRoleId_fkey" FOREIGN KEY ("userRoleId") REFERENCES public."m_userRole"(id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: gcasGroupSystemRelation gcasGroupSystemRelation_systemId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."gcasGroupSystemRelation"
    ADD CONSTRAINT "gcasGroupSystemRelation_systemId_fkey" FOREIGN KEY ("systemId") REFERENCES public.system(id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: gcasGroupSystemRelation gcasGroupSystemRelation_groupId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."gcasGroupSystemRelation"
    ADD CONSTRAINT "gcasGroupSystemRelation_groupId_fkey" FOREIGN KEY ("groupId") REFERENCES public."gcasGroup"(id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: project project_localGovernmentId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.project
    ADD CONSTRAINT "project_localGovernmentId_fkey" FOREIGN KEY ("localGovernmentId") REFERENCES public."m_localGovernment"(id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: projectCost projectCost_projectId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."projectCost"
    ADD CONSTRAINT "projectCost_projectId_fkey" FOREIGN KEY ("projectId") REFERENCES public.project(id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: projectSystemRelation projectSystemRelation_projectId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."projectSystemRelation"
    ADD CONSTRAINT "projectSystemRelation_projectId_fkey" FOREIGN KEY ("projectId") REFERENCES public.project(id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: projectSystemRelation projectSystemRelation_systemId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."projectSystemRelation"
    ADD CONSTRAINT "projectSystemRelation_systemId_fkey" FOREIGN KEY ("systemId") REFERENCES public.system(id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: system system_localGovernmentId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.system
    ADD CONSTRAINT "system_localGovernmentId_fkey" FOREIGN KEY ("localGovernmentId") REFERENCES public."m_localGovernment"(id) ON UPDATE CASCADE ON DELETE SET NULL;

--
-- Name: systemBasicInformation systemBasicInformation_projectId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."systemBasicInformation"
    ADD CONSTRAINT "systemBasicInformation_projectId_fkey" FOREIGN KEY ("projectId") REFERENCES public.project(id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- PostgreSQL database dump complete
--

