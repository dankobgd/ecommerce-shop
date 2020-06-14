CREATE TYPE gender AS ENUM ('m', 'f');

CREATE TABLE public.user (
  id int generated always as identity primary key,
  first_name varchar(255),
  last_name varchar(255),
  username varchar(255),
  email varchar(255) unique not null,
  password text not null,
  gender gender,
  locale varchar(5) not null,
  avatar_url text,
  active bool not null,
  email_verified bool not null,
  failed_attempts int,
  last_login_at timestamptz,
  created_at timestamptz not null,
  updated_at timestamptz not null,
  deleted_at timestamptz
);

CREATE TABLE public.token (
  id int generated always as identity primary key,
  user_id int not null,
  token text not null,
  type varchar(64) not null,
  created_at timestamptz not null,
  expires_at timestamptz not null,
  foreign key (user_id) references public.user (id) on delete cascade
);

CREATE TABLE public.role (
  id int generated always as identity primary key,
  name varchar(30)
);

CREATE TABLE public.user_role (
  user_id int not null,
  role_id int not null,
  primary key (user_id, role_id),
  foreign key (user_id) references public.user (id),
  foreign key (role_id) references public.role (id)
);

CREATE TABLE public.address (
  id int generated always as identity primary key,
  user_id int,
  country varchar(255),
  city varchar(255),
  address_1 varchar(255),
  address_2 varchar(255),
  zip varchar(30),
  longitude numeric(11, 8),
  latitude numeric(11, 8),
  foreign key (user_id) references public.user (id)
);

CREATE TABLE public.item (
  id int generated always as identity primary key,
  user_id int,
  type_id int,
  size_id int,
  color_id int,
  description text,
  sku varchar(30)
);

CREATE TABLE public.item_info (
  id int generated always as identity primary key,
  item_id int,
  price int,
  description text,
  foreign key (item_id) references public.item (id)
);

CREATE TABLE public.item_images (
  id int generated always as identity primary key,
  item_id int,
  url text,  
  foreign key (item_id) references public.item (id)
);

CREATE TABLE public.related_item (
  id int generated always as identity primary key,
  item_id int,
  related_item_id int,  
  foreign key (item_id) references public.item (id),
  foreign key (related_item_id) references public.item (id)
);

CREATE TABLE public.manufacturer (
  id int generated always as identity primary key,
  name varchar(100),
  type varchar(50),
  email text,
  website_url text,
  address text,
  description text
);

CREATE TABLE public.color (
  id int generated always as identity primary key
);

CREATE TABLE public.size (
  id int generated always as identity primary key
);
