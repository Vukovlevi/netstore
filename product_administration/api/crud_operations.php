<?php
header('Access-Control-Allow-Origin: *');
header('Access-Control-Allow-Headers: Content-Type, Authorization');
header('Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS');

if ($_SERVER['REQUEST_METHOD'] === 'OPTIONS') {
    http_response_code(200);
    exit();
}

require './sql_functions.php';
require './read_cookies.php';
require './crud/category.php';
require './crud/sub_category.php';
require './crud/product_type.php';
require './crud/storing_condition.php';
require './crud/brand.php';
require './crud/product.php';
require './crud/search_product.php';

$method = $_SERVER['REQUEST_METHOD'];
$uri = explode('/', parse_url($_SERVER['REQUEST_URI'], PHP_URL_PATH));
$body = json_decode(file_get_contents('php://input'), true);

//$isAuth = authentication();
//if(!$isAuth) return http_response_code(401);

switch(end($uri)) {
    case 'category':
        handleCategory($method, $body);
        break;
    case 'sub_category':
        handleSubCategory($method, $body);
        break;
    case 'product_type':
        handleProductType($method, $body);
        break;
    case 'storing_condition':
        handleStoringCondition($method,$body);
        break;
    case 'brand':
        handleBrand($method,$body);
        break;
    case 'product':
        handleProduct($method, $body);
        break;
    case 'search_product':
        handleSearchProduct($method, $body);
        break;
    case 'auth':
        return http_response_code(200);
    default:
        break;
}