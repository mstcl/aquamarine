package ansi

const Escape string = "\033["

const Regular string = "0;"
const Bold string = "1;"

const BoldYellow string = Escape + Bold + "33m"

const Red string = Escape + Regular + "31m"

const BoldBlue string = Escape + Bold + "34m"

const BrightWhite string = Escape + Regular + "97m"

const Reset string = Escape + Regular + "0m"
