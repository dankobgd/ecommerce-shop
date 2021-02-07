create type gender AS ENUM ('m', 'f');


create table public.category (
  id int generated always as identity primary key,
  name varchar(64) not null,
  slug varchar(64) not null,
  logo text not null,
  logo_public_id text not null,
  description text,
  is_featured bool default false not null,
  properties jsonb,
  created_at timestamptz not null,
  updated_at timestamptz not null,
  unique (name, slug)
);

create table public.brand (
  id int generated always as identity primary key,
  name varchar(64) not null,
  slug varchar(64) not null,
  type varchar(64) not null,
  description text,
  email text not null,
  website_url text not null,
  logo text not null,
  logo_public_id text not null,
  created_at timestamptz not null,
  updated_at timestamptz not null,
  unique (name, slug)
);

create table public.tag (
  id int generated always as identity primary key,
  name varchar(64) not null,
  slug varchar(64) not null,
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
  foreign key (user_id) references public.user (id) on delete cascade,
  foreign key (address_id) references public.address (id) on delete cascade
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
  name varchar(64) not null,  
  slug varchar(64) not null,
  image_url text not null,
  image_public_id text not null,
  description text,
  in_stock bool default true not null,
  sku text not null,
  is_featured bool default false not null,
  created_at timestamptz not null,
  updated_at timestamptz not null,
  properties jsonb,
  foreign key (brand_id) references public.brand (id) on delete cascade,
  foreign key (category_id) references public.category (id) on delete cascade
);

create table public.product_pricing (
  id int generated always as identity primary key,
  product_id int not null,
  price int not null,
  original_price int not null,
  sale_starts timestamptz not null,
  sale_ends timestamptz not null,
  foreign key (product_id) references public.product (id) on delete cascade
);

create table public.product_tag (
  product_id int not null,
  tag_id int not null,  
  foreign key (product_id) references public.product (id) on delete cascade,
  foreign key (tag_id) references public.tag (id) on delete cascade,
  unique (product_id, tag_id)
);

create table public.product_image (
  id int generated always as identity primary key,
  product_id int not null,
  url text not null,
  public_id text not null,
  created_at timestamptz not null,
  updated_at timestamptz not null,
  foreign key (product_id) references public.product (id) on delete cascade
);

create table public.product_review (
  id int generated always as identity primary key,
  user_id int not null,
  product_id int not null,
  rating int not null,
  title text not null,
  comment text not null,
  created_at timestamptz not null,
  updated_at timestamptz not null,
  foreign key (user_id) references public.user (id) on delete cascade,
  foreign key (product_id) references public.product (id) on delete cascade,
  check (rating >= 1 and rating <= 5),
  unique (user_id, product_id)
);

create table public.product_wishlist (
  id int generated always as identity primary key,
  user_id int not null,
  product_id int not null,
  foreign key (user_id) references public.user (id) on delete cascade,
  foreign key (product_id) references public.product (id) on delete cascade,
  unique (user_id, product_id)
);

create table public.promotion (
  promo_code varchar(30) primary key,
  type varchar(30) not null,
  amount int not null,  
  description text,
  starts_at timestamptz not null,
  ends_at timestamptz not null
);

create table public.promotion_detail (
  user_id int,
  promo_code varchar(30),
  foreign key (user_id) references public.user (id) on delete cascade,
  foreign key (promo_code) references public.promotion (promo_code) on delete cascade,
  unique (user_id, promo_code)  
);

create table public.order (
  id int generated always as identity primary key,
  user_id int not null,
  promo_code varchar(30),
  promo_code_type varchar(30),
  promo_code_amount int,  
  status varchar(30) default 'pending' not null,
  subtotal int not null,
  total int not null,
  shipped_at timestamptz,
  created_at timestamptz not null,
  payment_method_id text not null,
  payment_intent_id text not null,
  receipt_url text not null,
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
  foreign key (user_id) references public.user (id),
  foreign key (promo_code) references public.promotion (promo_code)
);

create table public.order_detail (
  order_id int,
  product_id int,
  quantity int not null,
  history_price int not null,
  history_sku text not null,
  primary key (order_id, product_id)
);


-- ideally use materialized view and refresh it but whatever
create view product_search_view as
select
p.*,
b.name AS brand_name,
b.slug AS brand_slug,
b.type AS brand_type,
b.description AS brand_description,
b.email AS brand_email,
b.logo AS brand_logo,
b.website_url AS brand_website_url,
b.created_at AS brand_created_at,
b.updated_at AS brand_updated_at,
c.name AS category_name,
c.slug AS category_slug,
c.description AS category_description,
c.logo AS category_logo,
c.created_at AS category_created_at,
c.updated_at AS category_updated_at,
pp.id AS pricing_id,
pp.product_id AS pricing_product_id,
pp.price AS pricing_price,
pp.original_price AS pricing_original_price,
pp.sale_starts AS pricing_sale_starts,
pp.sale_ends AS pricing_sale_ends,
(
setweight(to_tsvector(coalesce(p.name,'')), 'A') ||
setweight(to_tsvector(coalesce(p.description,'')), 'B') ||
setweight(to_tsvector(coalesce(c.name,'')), 'C') ||
setweight(to_tsvector(coalesce(b.name,'')), 'D')
) as tsv
FROM product p
LEFT JOIN product_pricing pp ON p.id = pp.product_id
LEFT JOIN brand b ON p.brand_id = b.id
LEFT JOIN category c ON p.category_id = c.id
WHERE CURRENT_TIMESTAMP BETWEEN pp.sale_starts AND pp.sale_ends;
