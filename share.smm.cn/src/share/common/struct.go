package share

import (

)

type Stringer interface{
    ToString()string
}

type Totaler interface{
    GetTotal()int
}
