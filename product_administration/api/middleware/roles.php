<?php
define('ROLE_UZLETVEZETO', 'Üzletvezető');
define('ROLE_HR', 'HR-es');
define('ROLE_RAKTARVEZETO', 'Raktárvezető');
define('ROLE_RAKTARKEZELO', 'Raktárkezelő');
define('ROLE_PENZTAROS', 'Pénztáros');
define('ROLE_EGYEB', 'Egyéb dolgozó');

function getCurrentUserRole() {
    return $_REQUEST['user']['name'] ?? null;
}

function requireRole($allowedRoles) {
    if (!isset($_REQUEST['user'])) {
        http_response_code(401);
        echo json_encode(['message' => 'Nincs bejelentkezve!'], JSON_UNESCAPED_UNICODE);
        return false;
    }

    $userRole = getCurrentUserRole();

    if (!in_array($userRole, $allowedRoles)) {
        http_response_code(403);
        echo json_encode(['message' => 'Nincs jogosultsága ehhez a művelethez!'], JSON_UNESCAPED_UNICODE);
        return false;
    }

    return true;
}

function requireAdmin() {
    return requireRole([ROLE_UZLETVEZETO]);
}

function requireHRAccess() {
    return requireRole([ROLE_UZLETVEZETO, ROLE_HR]);
}

function requireWarehouseManagerAccess() {
    return requireRole([ROLE_UZLETVEZETO, ROLE_RAKTARVEZETO]);
}

function requireProductManagementAccess() {
    return requireRole([ROLE_UZLETVEZETO, ROLE_RAKTARVEZETO, ROLE_RAKTARKEZELO]);
}

function requireSearchAccess() {
    return requireRole([ROLE_UZLETVEZETO, ROLE_HR, ROLE_RAKTARVEZETO, ROLE_RAKTARKEZELO, ROLE_PENZTAROS, ROLE_EGYEB]);
}

function requireQuantityDecreaseAccess() {
    return requireRole([ROLE_UZLETVEZETO, ROLE_RAKTARVEZETO, ROLE_RAKTARKEZELO, ROLE_PENZTAROS]);
}

function checkResourceAccess($resource, $method) {
    if ($method === 'GET') {
        return requireSearchAccess();
    }

    switch ($resource) {
        case 'category':
        case 'sub_category':
        case 'product_type':
        case 'storing_condition':
        case 'brand':
        case 'ingredient':
            return requireWarehouseManagerAccess();

        case 'product':
            return requireProductManagementAccess();

        case 'user':
        case 'contract':
        case 'contract_type':
            return requireHRAccess();

        case 'store_detail':
        case 'open_hour':
            return requireAdmin();

        default:
            return requireAdmin();
    }
}

function canWrite($resource) {
    $userRole = getCurrentUserRole();

    if ($userRole === ROLE_UZLETVEZETO) {
        return true;
    }

    switch ($resource) {
        case 'category':
        case 'sub_category':
        case 'product_type':
        case 'storing_condition':
        case 'brand':
        case 'ingredient':
            return $userRole === ROLE_RAKTARVEZETO;

        case 'product':
            return in_array($userRole, [ROLE_RAKTARVEZETO, ROLE_RAKTARKEZELO]);

        case 'user':
        case 'contract':
        case 'contract_type':
            return $userRole === ROLE_HR;

        default:
            return false;
    }
}

function canDecreaseQuantity() {
    $userRole = getCurrentUserRole();
    return in_array($userRole, [ROLE_UZLETVEZETO, ROLE_RAKTARVEZETO, ROLE_RAKTARKEZELO, ROLE_PENZTAROS]);
}
