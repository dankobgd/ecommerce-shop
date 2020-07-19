-- custom types
drop type gender;
create type gender AS ENUM ('m', 'f');


-- valid possible reference table values
create table public.category (
  id int generated always as identity primary key,
  name varchar(50) unique
);
create table public.color (
  id int generated always as identity primary key,
  value varchar(50) unique
);
create table public.size (
  id int generated always as identity primary key,
  value varchar(20) unique
);
create table public.shoe_size (
  id int generated always as identity primary key,
  value varchar(20) unique
);
create table public.jeans_size (
  id int generated always as identity primary key,
  value varchar(20) unique
);


-- app tables
create table public.address (
  id int generated always as identity primary key,
  line_1 text not null,
  line_2 text,
  city text not null,
  country text not null,
  state text,
  zip text,
  latitude numeric(11, 8),
  longitude numeric(11, 8),
  phone text,
  created_at timestamptz not null,
  updated_at timestamptz not null,
  deleted_at timestamptz
);

create table public.address_type (
  id int generated always as identity primary key,
  address_id int not null,
  name varchar(30) not null,
  foreign key (address_id) references public.address (id) on delete cascade
);

create table public.user (
  id int generated always as identity primary key,
  first_name varchar(255),
  last_name varchar(255),
  username varchar(255),
  email varchar(255) unique not null,
  password text not null,
  role varchar(20) default 'user' not null,
  gender gender,
  locale varchar(5) default 'en' not null,
  avatar_url text,
  active bool not null,
  email_verified bool default false not null,
  failed_attempts int default 0 not null,
  last_login_at timestamptz not null,
  created_at timestamptz not null,
  updated_at timestamptz not null,
  deleted_at timestamptz
);

create table public.user_address (
  user_id int not null,
  address_id int not null,
  address_type_id int not null,
  foreign key (user_id) references public.user (id) on delete cascade,
  foreign key (address_id) references public.address (id) on delete cascade,
  foreign key (address_type_id) references public.address_type (id) on delete cascade
);

create table public.token (
  id int generated always as identity primary key,
  user_id int not null,
  token text not null,
  type varchar(64) not null,
  created_at timestamptz not null,
  expires_at timestamptz not null,
  foreign key (user_id) references public.user (id) on delete cascade
);

create table public.product (
  id int generated always as identity primary key,
  name varchar(255) not null,
  slug varchar(50),
  image_url text not null,
  description text,
  price int not null,
  stock int default 0 not null,
  sku text not null,
  is_featured bool default false not null,
  created_at timestamptz not null,
  updated_at timestamptz not null
);

create table public.product_info (
  id int generated always as identity primary key,
  product_id int not null,
  info text not null,
  foreign key (product_id) references public.product (id)
);

create table public.product_category (
  id int generated always as identity primary key,
  product_id int not null,
  name varchar(50),
  slug varchar(50),
  description text,
  foreign key (name) references public.category (name),
  foreign key (product_id) references public.product (id) on delete cascade
);

create table public.product_tag (
  id int generated always as identity primary key,
  product_id int not null,
  name varchar(255) not null,
  created_at timestamptz not null,
  updated_at timestamptz not null,
  foreign key (product_id) references public.product (id) on delete cascade
);

create table public.product_brand (
  id int generated always as identity primary key,
  product_id int not null,
  name varchar(255) not null,
  slug varchar(50),
  type varchar(50) not null,
  description text,
  email text not null,
  website_url text not null,
  created_at timestamptz not null,
  updated_at timestamptz not null,
  foreign key (product_id) references public.product (id) on delete cascade
);

create table public.product_image (
  id int generated always as identity primary key,
  product_id int not null,
  url text not null,
  created_at timestamptz not null,
  updated_at timestamptz not null,
  foreign key (product_id) references public.product (id) on delete cascade
);

create table public.product_discount (
  id int generated always as identity primary key,
  product_id int not null,
  fixed_value int,
  percentage_value int,
  description text,
  starts_at timestamptz,
  ends_at timestamptz,
  foreign key (product_id) references public.product (id) on delete cascade
);

create table public.related_product (
  id int generated always as identity primary key,
  product_id int not null,
  related_product_id int not null,
  foreign key (product_id) references public.product (id),
  foreign key (related_product_id) references public.product (id)
);

create table public.order (
  id int generated always as identity primary key,
  user_id int not null,  
  status varchar(30) default 'pending' not null,
  total int,
  shipped_at timestamptz,
  created_at timestamptz not null,
  billing_address_line_1 text not null,
  billing_address_line_2 text,
  billing_address_city text not null,
  billing_address_country text not null,
  billing_address_state text,
  billing_address_zip text,
  billing_address_latitude numeric(11, 8),
  billing_address_longitude numeric(11, 8),  
  shipping_address_line_1 text not null,
  shipping_address_line_2 text,
  shipping_address_city text not null,
  shipping_address_country text not null,
  shipping_address_state text,
  shipping_address_zip text,
  shipping_address_latitude numeric(11, 8),
  shipping_address_longitude numeric(11, 8),  
  foreign key (user_id) references public.user (id)
);

create table public.order_detail (
  order_id int,
  product_id int,
  quantity int not null,
  original_price int not null,
  original_sku text not null,
  constraint pk_order_detail primary key(order_id, product_id)
);


-- populate tables
COPY public.category (name) FROM '/datasource/category.csv';
COPY public.size (value) FROM '/datasource/size.csv';
COPY public.shoe_size (value) FROM '/datasource/shoe_size.csv';
COPY public.jeans_size (value) FROM '/datasource/jeans_size.csv';
COPY public.color (value) FROM '/datasource/color.csv';
