<?php

require __DIR__ . '/../vendor/autoload.php';

require_once __DIR__ . '/support/DbStub.php';
require_once __DIR__ . '/support/HandlerTestCase.php';

if (!function_exists('getData')) {
    function getData($operation, $types = null, $data = null) {
        DbStub::$getDataCalls[] = [
            'sql'   => $operation,
            'types' => $types,
            'data'  => $data,
        ];
        if (empty(DbStub::$getDataResponses)) {
            return [];
        }
        return array_shift(DbStub::$getDataResponses);
    }
}

if (!function_exists('changeData')) {
    function changeData($operation, $types = null, $data = null) {
        DbStub::$changeDataCalls[] = [
            'sql'   => $operation,
            'types' => $types,
            'data'  => $data,
        ];
        if (empty(DbStub::$changeDataResponses)) {
            return true;
        }
        return array_shift(DbStub::$changeDataResponses);
    }
}

require_once __DIR__ . '/../api/middleware/roles.php';
require_once __DIR__ . '/../api/crud/category.php';
require_once __DIR__ . '/../api/crud/sub_category.php';
require_once __DIR__ . '/../api/crud/product_type.php';
require_once __DIR__ . '/../api/crud/storing_condition.php';
require_once __DIR__ . '/../api/crud/brand.php';
require_once __DIR__ . '/../api/crud/product.php';
require_once __DIR__ . '/../api/crud/search_product.php';
