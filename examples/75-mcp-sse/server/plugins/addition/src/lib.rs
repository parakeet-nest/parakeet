use extism_pdk::*;
use serde::{Deserialize, Serialize};


// use json data for inputs and outputs
#[derive(FromBytes, Deserialize, PartialEq, Debug)]
#[encoding(Json)]
struct Add {
    left: i32,
    right: i32,
}
#[derive(ToBytes, Serialize, PartialEq, Debug)]
#[encoding(Json)]
struct Sum {
    value: i32,
}

#[plugin_fn]
pub fn add(input: Add) -> FnResult<Sum> {
    Ok(Sum {
        value: input.left + input.right,
    })
}
