<?php

final class ProductTest extends HandlerTestCase
{
    private function validPostBody(array $overrides = []): array
    {
        return array_merge([
            'name'        => 'Teszt Termék',
            'description' => 'Leírás',
            'amount'      => 10,
            'size'        => '1',
            'size_type'   => 'l',
            'price'       => 500,
            'discount'    => 0,
            'type_id'     => 2,
            'brand_id'    => 3,
        ], $overrides);
    }

    public function testPostMissingRequiredFieldsReturns400(): void
    {
        $res = $this->runAndCapture(fn () => handleProduct('POST', ['name' => 'x']));

        $this->assertSame(400, $res['code']);
        $this->assertStringContainsString('Hiányzó adat', $res['json']['message']);
        $this->assertSame([], DbStub::$getDataCalls);
    }

    public function testPostDuplicateActiveProductReturns409(): void
    {
        DbStub::queueGetData([
            ['id' => 5, 'deleted_at' => null],
        ]);

        $res = $this->runAndCapture(fn () => handleProduct('POST', $this->validPostBody()));

        $this->assertSame(409, $res['code']);
        $this->assertStringContainsString('aktív termék', $res['json']['message']);
    }

    public function testPostDuplicateSoftDeletedProductReturns409(): void
    {
        DbStub::queueGetData([
            ['id' => 5, 'deleted_at' => '2026-01-15'],
        ]);

        $res = $this->runAndCapture(fn () => handleProduct('POST', $this->validPostBody()));

        $this->assertSame(409, $res['code']);
        $this->assertStringContainsString('törölve', $res['json']['message']);
    }

    public function testPostInsertsAndReturnsNewProduct(): void
    {
        DbStub::queueGetData(
            [],
            [['id' => 42, 'name' => 'Teszt Termék']]
        );
        DbStub::queueChangeData(true);

        $res = $this->runAndCapture(fn () => handleProduct(
            'POST',
            $this->validPostBody(['discount' => 20, 'warranty' => '2027', 'expires_at' => '2026-12-31T08:00:00'])
        ));

        $this->assertSame(201, $res['code']);
        $this->assertSame(['id' => 42, 'name' => 'Teszt Termék'], $res['json']);

        $insert = DbStub::$changeDataCalls[0];
        $this->assertStringContainsString('INSERT INTO product', $insert['sql']);
        $this->assertSame('ssisssidsii', $insert['types']);

        $this->assertEqualsWithDelta(0.2, $insert['data'][7], 0.0001);

        $this->assertSame('2027-01-01', $insert['data'][8]);

        $this->assertSame('2026-12-31', $insert['data'][5]);
    }

    public function testPostCoercesInvalidDateToNull(): void
    {
        DbStub::queueGetData([], [['id' => 1, 'name' => 'x']]);
        DbStub::queueChangeData(true);

        $this->runAndCapture(fn () => handleProduct(
            'POST',
            $this->validPostBody(['name' => 'x', 'warranty' => 'not-a-date', 'expires_at' => 'undefined'])
        ));

        $insert = DbStub::$changeDataCalls[0];
        $this->assertNull($insert['data'][5]);
        $this->assertNull($insert['data'][8]);
    }

    public function testPutReturns404WhenProductMissing(): void
    {
        DbStub::queueGetData([]);

        $res = $this->runAndCapture(fn () => handleProduct(
            'PUT',
            $this->validPostBody(['id' => 99])
        ));

        $this->assertSame(404, $res['code']);
        $this->assertStringContainsString('Nincs ilyen', $res['json']['message']);
    }

    public function testPutReturns409WhenNameIsTaken(): void
    {
        DbStub::queueGetData(
            [['id' => 7]],
            [['id' => 8]]
        );

        $res = $this->runAndCapture(fn () => handleProduct(
            'PUT',
            $this->validPostBody(['id' => 7])
        ));

        $this->assertSame(409, $res['code']);
        $this->assertStringContainsString('foglalt', $res['json']['message']);
    }

    public function testPutUpdatesProductSuccessfully(): void
    {
        DbStub::queueGetData([['id' => 7]], []);
        DbStub::queueChangeData(true);

        $res = $this->runAndCapture(fn () => handleProduct(
            'PUT',
            $this->validPostBody(['id' => 7, 'name' => 'Frissített'])
        ));

        $this->assertSame(200, $res['code']);
        $this->assertSame('Sikeres frissítés', $res['json']['message']);

        $update = DbStub::$changeDataCalls[0];
        $this->assertStringStartsWith('UPDATE product', $update['sql']);
        $this->assertSame(7, end($update['data']));
    }

    public function testDeleteWithoutIdReturns400(): void
    {
        $res = $this->runAndCapture(fn () => handleProduct('DELETE', null));

        $this->assertSame(400, $res['code']);
        $this->assertSame([], DbStub::$getDataCalls);
    }

    public function testDeleteNonExistentReturns404(): void
    {
        $_GET['id'] = 123;
        DbStub::queueGetData([]);

        $res = $this->runAndCapture(fn () => handleProduct('DELETE', null));

        $this->assertSame(404, $res['code']);
    }

    public function testDeleteSoftDeletes(): void
    {
        $_GET['id'] = 10;
        DbStub::queueGetData([['id' => 10]]);
        DbStub::queueChangeData(true);

        $res = $this->runAndCapture(fn () => handleProduct('DELETE', null));

        $this->assertSame(204, $res['code']);
        $update = DbStub::$changeDataCalls[0];
        $this->assertStringContainsString('SET deleted_at = CURDATE()', $update['sql']);
        $this->assertSame([10], $update['data']);
    }

    public function testGetByIdJoinsTypeAndBrand(): void
    {
        $_GET['id'] = 1;
        DbStub::queueGetData([
            ['id' => 1, 'name' => 'P', 'type_name' => 'T', 'brand_name' => 'B'],
        ]);

        $res = $this->runAndCapture(fn () => handleProduct('GET', null));

        $this->assertSame(200, $res['code']);
        $this->assertSame('T', $res['json']['type_name']);
        $this->assertSame('B', $res['json']['brand_name']);

        $sql = DbStub::$getDataCalls[0]['sql'];
        $this->assertStringContainsString('LEFT JOIN product_type', $sql);
        $this->assertStringContainsString('LEFT JOIN brand', $sql);
        $this->assertStringContainsString('deleted_at IS NULL', $sql);
    }

    public function testGetByIdNotFoundReturns404(): void
    {
        $_GET['id'] = 999;
        DbStub::queueGetData([]);

        $res = $this->runAndCapture(fn () => handleProduct('GET', null));

        $this->assertSame(404, $res['code']);
    }

    public function testUnsupportedMethodReturns405(): void
    {
        $res = $this->runAndCapture(fn () => handleProduct('PATCH', null));

        $this->assertSame(405, $res['code']);
    }
}
