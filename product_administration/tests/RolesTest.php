<?php

use PHPUnit\Framework\Attributes\DataProvider;

final class RolesTest extends HandlerTestCase
{
    private function loginAs(string $role): void
    {
        $_REQUEST['user'] = ['name' => $role, 'id' => 1, 'username' => 'tester'];
    }

    public function testRequireRoleReturnsFalseWhenUnauthenticated(): void
    {
        $res = $this->runAndCapture(function () {
            $this->assertFalse(requireRole([ROLE_UZLETVEZETO]));
        });

        $this->assertSame(401, $res['code']);
        $this->assertStringContainsString('Nincs bejelentkezve', $res['json']['message']);
    }

    public function testRequireRoleReturnsFalseForDisallowedRole(): void
    {
        $this->loginAs(ROLE_PENZTAROS);

        $res = $this->runAndCapture(function () {
            $this->assertFalse(requireRole([ROLE_UZLETVEZETO, ROLE_HR]));
        });

        $this->assertSame(403, $res['code']);
        $this->assertStringContainsString('Nincs jogosultsága', $res['json']['message']);
    }

    public function testRequireRoleSucceedsForAllowedRole(): void
    {
        $this->loginAs(ROLE_HR);

        $res = $this->runAndCapture(function () {
            $this->assertTrue(requireRole([ROLE_UZLETVEZETO, ROLE_HR]));
        });

        $this->assertSame(200, $res['code']);
        $this->assertSame('', $res['body'], 'No body should be written on success');
    }

    public function testCheckResourceAccessGetIsOpenToAllStaff(): void
    {
        $this->loginAs(ROLE_EGYEB);

        $this->runAndCapture(function () {
            $this->assertTrue(checkResourceAccess('product', 'GET'));
        });
    }

    public function testCheckResourceAccessProductWriteRequiresWarehouseRole(): void
    {
        $this->loginAs(ROLE_PENZTAROS);
        $res = $this->runAndCapture(function () {
            $this->assertFalse(checkResourceAccess('product', 'POST'));
        });
        $this->assertSame(403, $res['code']);

        http_response_code(200);

        $this->loginAs(ROLE_RAKTARKEZELO);
        $this->runAndCapture(function () {
            $this->assertTrue(checkResourceAccess('product', 'POST'));
        });
    }

    public function testCheckResourceAccessCategoryWriteRequiresWarehouseManagerOrAdmin(): void
    {
        $this->loginAs(ROLE_RAKTARKEZELO);
        $res = $this->runAndCapture(function () {
            $this->assertFalse(checkResourceAccess('category', 'PUT'));
        });
        $this->assertSame(403, $res['code']);

        http_response_code(200);

        $this->loginAs(ROLE_RAKTARVEZETO);
        $this->runAndCapture(function () {
            $this->assertTrue(checkResourceAccess('category', 'PUT'));
        });
    }

    public function testCheckResourceAccessUserResourceRequiresHR(): void
    {
        $this->loginAs(ROLE_RAKTARVEZETO);
        $this->runAndCapture(function () {
            $this->assertFalse(checkResourceAccess('user', 'POST'));
        });

        $this->loginAs(ROLE_HR);
        $this->runAndCapture(function () {
            $this->assertTrue(checkResourceAccess('user', 'POST'));
        });
    }

    public function testCheckResourceAccessUnknownResourceRequiresAdmin(): void
    {
        $this->loginAs(ROLE_HR);
        $this->runAndCapture(function () {
            $this->assertFalse(checkResourceAccess('something_unknown', 'POST'));
        });

        $this->loginAs(ROLE_UZLETVEZETO);
        $this->runAndCapture(function () {
            $this->assertTrue(checkResourceAccess('something_unknown', 'POST'));
        });
    }

    public static function canWriteMatrix(): array
    {
        return [
            'Admin can write products'               => [ROLE_UZLETVEZETO, 'product', true],
            'Admin can write users'                  => [ROLE_UZLETVEZETO, 'user', true],
            'WarehouseManager writes categories'     => [ROLE_RAKTARVEZETO, 'category', true],
            'WarehouseManager cannot write users'    => [ROLE_RAKTARVEZETO, 'user', false],
            'WarehouseKeeper writes products'        => [ROLE_RAKTARKEZELO, 'product', true],
            'WarehouseKeeper cannot write brand'     => [ROLE_RAKTARKEZELO, 'brand', false],
            'HR writes users'                        => [ROLE_HR, 'user', true],
            'HR cannot write products'               => [ROLE_HR, 'product', false],
            'Cashier cannot write anything'          => [ROLE_PENZTAROS, 'product', false],
            'Cashier cannot write category'          => [ROLE_PENZTAROS, 'category', false],
            'Egyeb cannot write anything'            => [ROLE_EGYEB, 'product', false],
        ];
    }

    #[DataProvider('canWriteMatrix')]
    public function testCanWriteMatrix(string $role, string $resource, bool $expected): void
    {
        $this->loginAs($role);
        $this->assertSame($expected, canWrite($resource));
    }

    public function testCanDecreaseQuantityOnlyForOperationalRoles(): void
    {
        $allowed = [ROLE_UZLETVEZETO, ROLE_RAKTARVEZETO, ROLE_RAKTARKEZELO, ROLE_PENZTAROS];
        $denied  = [ROLE_HR, ROLE_EGYEB];

        foreach ($allowed as $role) {
            $this->loginAs($role);
            $this->assertTrue(canDecreaseQuantity(), "$role should be able to decrease quantity");
        }
        foreach ($denied as $role) {
            $this->loginAs($role);
            $this->assertFalse(canDecreaseQuantity(), "$role should NOT be able to decrease quantity");
        }
    }
}
