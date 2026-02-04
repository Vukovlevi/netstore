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
require './middleware/roles.php';
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
$resource = end($uri);

// Authentication check
$isAuth = authentication();
if (!$isAuth && $resource !== 'auth') {
    http_response_code(401);
    echo json_encode(['message' => 'Nincs bejelentkezve!'], JSON_UNESCAPED_UNICODE);
    exit();
}

// Role-based access control for write operations
if ($method !== 'GET' && $resource !== 'auth') {
    if (!checkResourceAccess($resource, $method)) {
        exit();
    }
}

switch($resource) {
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
        http_response_code(200);
        if ($isAuth) {
            echo json_encode($_REQUEST['user'], JSON_UNESCAPED_UNICODE);
        }
        break;
    default:
        http_response_code(404);
        break;
}