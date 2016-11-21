package trade

const(
    DefaultPageSize = 20
    MaxPageView     = 10
)

type Page struct{
    First       int
    Last        int
    Prev        int
    Next        int
    From        int
    To          int
    Pages       []int
    Now         int
    Total       int
    Size        int
    UrlLeft     string
    UrlRight    string
}

func (p *Page) DefaultSize () int{
    return DefaultPageSize
}

func (p *Page) DefaultMax () int{
    return MaxPageView
}

func (p *Page) Set (now , total , max , size int , url_l , url_r string){
    p.Size = size
    p.UrlLeft = url_l
    p.UrlRight = url_r
    p.Now = now

    p.First = 1

    p.Last  = ((total - 1) / size) + 1

    p.Total = total

    if p.Now == 1 {

        p.Prev  =  1
    }else{
        p.Prev  =  p.Now - 1
    }

    if p.Now == p.Last {

        p.Next  = p.Last
    }else{
        p.Next  = p.Now + 1
    }

    if p.Now <= max / 2 {

        p.From = 1

        if p.Last > max {

            p.To   =  max
        }else{
            p.To   = p.Last
        }
    }else if p.Now >= p.Last - max / 2 {

        p.To   = p.Last

        if p.Last > max {

            p.From =  p.Last - max
        }else{
            p.From = 1
        }

    }else{
        p.From = p.Now - max / 2
        p.To   = p.Now + max / 2
    }

    p.Pages = make( []int , 0 , max )

    for i := p.From; i <= p.To; i++ {
        p.Pages = append(p.Pages , i )
    }

}