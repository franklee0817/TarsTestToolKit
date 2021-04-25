<?php

namespace php\impl;

use php\servant\TestUnits\php\testObj\testServant;

class test implements testServant
{
    public function ping($req)
    {
        return $req;
    }
}
