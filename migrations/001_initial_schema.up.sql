create type gender AS ENUM ('m', 'f');

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

insert into public.category (name) values ('shirt'), ('tshirt'), ('shorts'), ('jeans'), ('shoes'), ('jacket');
insert into public.size (value) values ('xs'), ('s'), ('sm'), ('m'), ('md'), ('l'), ('xl'), ('2xl'), ('3xl'), ('4xl'), ('5xl');
insert into public.shoe_size (value) values ('20'), ('21'), ('22'), ('23'), ('24'), ('25'), ('26'), ('27'), ('28'), ('29'), ('30'), ('31'), ('32'), ('33'), ('34'), ('35'), ('36'), ('37'), ('38'), ('39'), ('40'), ('41'), ('42'), ('43'), ('44'), ('45'), ('46'), ('47'), ('48'), ('49'), ('50');
insert into public.jeans_size (value) values ('28/28'), ('30/28'), ('32/28'), ('34/28'), ('36/28'), ('38/28'), ('40/28'), ('30/30'), ('32/30'), ('34/30'), ('36/30'), ('38/30'), ('40/30'), ('32/32'), ('34/32'), ('36/32'), ('38/32'), ('40/32'), ('34/34'), ('36/36'), ('38/36'), ('40/36');
insert into public.color (value) values ('white'), ('black'), ('grey'), ('brown'), ('yellow'), ('red'), ('green'), ('blue'), ('orange'), ('purple'), ('pink');


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
  last_login_at timestamptz,
  created_at timestamptz not null,
  updated_at timestamptz not null,
  deleted_at timestamptz
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

create table public.contact (
  id int generated always as identity primary key,
  user_id int not null,
  country varchar(255) not null,
  city varchar(255) not null,
  address_1 varchar(255) not null,
  address_2 varchar(255),
  zip varchar(30),
  longitude numeric(11, 8),
  latitude numeric(11, 8),
  phone varchar(30),
  foreign key (user_id) references public.user (id)
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
  updated_at timestamptz not null,
  deleted_at timestamptz
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
  foreign key (product_id) references public.product (id),
  foreign key (name) references public.category (name) on delete cascade
);

create table public.product_tag (
  id int generated always as identity primary key,
  product_id int not null,
  name varchar(255) not null,
  created_at timestamptz not null,
  updated_at timestamptz not null,
  foreign key (product_id) references public.product (id)
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
  foreign key (product_id) references public.product (id)
);

create table public.product_image (
  id int generated always as identity primary key,
  product_id int not null,
  url text not null,
  foreign key (product_id) references public.product (id)
);

create table public.product_discount (
  id int generated always as identity primary key,
  product_id int not null,
  fixed_value int,
  percentage_value int,
  description text,
  starts_at timestamptz,
  ends_at timestamptz,
  foreign key (product_id) references public.product (id)
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
  shipping_address text not null,
  created_at timestamptz not null,
  foreign key (user_id) references public.user (id)
);

create table public.order_detail (
  order_id int,
  product_id int,
  quantity int not null,
  subtotal int not null,
  total int not null,
  provider varchar(50) not null,
  charge_ammount int not null,
  stripe_token text not null,
  stripe_token_type text not null,
  receipt_url text not null,
  constraint pk_order_detail primary key(order_id, product_id)
);

create table public.invoice (
  id int generated always as identity primary key,
  order_id int not null,
  code text not null,
  foreign key (order_id) references public.order (id)
);
