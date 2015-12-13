exception CalculationError {
  1: string what
}

service calculus {
    double calculate(1:string expr) throws (1:CalculationError err)
}
