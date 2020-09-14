create type gender AS ENUM ('m', 'f');

create table public.category (
  id int generated always as identity primary key,
  name varchar(50),
  slug varchar(50),
  logo text,
  description text,
  is_featured bool default false not null,
  created_at timestamptz not null,
  updated_at timestamptz not null
);

create table public.brand (
  id int generated always as identity primary key,
  name varchar(255) not null,
  slug varchar(50),
  type varchar(50) not null,
  description text,
  email text not null,
  website_url text not null,
  logo text,
  created_at timestamptz not null,
  updated_at timestamptz not null
);

create table public.tag (
  id int generated always as identity primary key,
  name varchar(50),
  slug varchar(50),
  description text,
  created_at timestamptz not null,
  updated_at timestamptz not null
);

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
  avatar_public_id text,
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
  brand_id int not null,
  category_id int not null,
  name varchar(255) not null,
  slug varchar(50),
  image_url text not null,
  description text,
  price int not null,
  in_stock bool default true not null,
  sku text not null,
  is_featured bool default false not null,
  created_at timestamptz not null,
  updated_at timestamptz not null,
  properties jsonb,
  foreign key (brand_id) references public.brand (id) on delete cascade,
  foreign key (category_id) references public.category (id) on delete cascade
);

create table public.product_info (
  id int generated always as identity primary key,
  product_id int not null,
  info text not null,
  foreign key (product_id) references public.product (id)
);

create table public.product_tag (
  id int generated always as identity primary key,
  product_id int not null,
  tag_id int not null,  
  foreign key (product_id) references public.product (id) on delete cascade,
  foreign key (tag_id) references public.tag (id) on delete cascade
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
  primary key(order_id, product_id)
);
