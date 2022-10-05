create table if not exists contacts (
    id integer primary key autoincrement, 
    fname text not null,
    lname text not null,
    phone text not null,
    email text not null,
    birthday text not null,
    address text not null,
    city text not null,
    state text not null,
    zipcode text not null,
    notes text not null
);

