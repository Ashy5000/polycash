// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
use crate::send::send;
use iced::widget::{button, column, container, row, rule, scrollable, text, text_input, Space};
use iced::{Alignment, Element, Sandbox, Settings};
use clipboard::{ClipboardContext, ClipboardProvider};

mod send;
mod sync;
mod key;
mod shorten;

pub fn main() -> iced::Result {
    let public_key = key::get_public_key();
    match public_key {
        Ok(_) => {}
        Err(e) => {
            println!("Key not found. Would you like to generate a new key? (y/n)");
            let mut input = String::new();
            std::io::stdin().read_line(&mut input).unwrap();
            if input.trim() == "y" {
                key::generate_key();
                println!("Continuing to GUI...")
            } else {
                println!("Exiting...");
                std::process::exit(0);
            }
        }
    }
    App::run(Settings::default())
}

struct App {
    balance: f64,
    amount: f64,
    address: String,
}

#[derive(Debug, Clone)]
enum Message {
    Sync,
    AmountChanged(String),
    AddressChanged(String),
    Send,
    Copy,
}

impl Sandbox for App {
    type Message = Message;

    fn new() -> Self {
        Self {
            balance: -1.0,
            amount: 0.0,
            address: "".to_string(),
        }
    }

    fn title(&self) -> String {
        String::from("GUI Wallet")
    }

    fn update(&mut self, message: Message) {
        match message {
            Message::Sync => {
                self.balance = sync::sync();
            }
            Message::AmountChanged(new_amount) => {
                let parsed_amount = new_amount.parse::<f64>();
                self.amount = match parsed_amount {
                    Ok(amount) => amount,
                    Err(_) => {
                        if new_amount == "" {
                            0.0
                        } else {
                            self.amount
                        }
                    }
                }
            }
            Message::AddressChanged(new_address) => {
                self.address = new_address;
            }
            Message::Send => {
                send(self.amount.clone(), self.address.clone());
            }
            Message::Copy => {
                let public_key = key::get_public_key().expect("Failed to get public key");
                let mut ctx: ClipboardContext = ClipboardProvider::new().unwrap();
                ctx.set_contents(public_key.clone()).unwrap();
                assert_eq!(ctx.get_contents().unwrap(), public_key);
            }
        }
    }

    fn view(&self) -> Element<Message> {
        scrollable(
            column![
                text("Balance: ").size(15),
                container(
                    row![
                        text({
                            if self.balance != -1.0 {
                                format!("{:.0}", self.balance)
                            } else {
                                "?".to_string()
                            }
                        })
                        .size(50),
                        text({
                            if self.balance != -1.0 {
                                format!(
                                    ".{}",
                                    format!("{:.2}", self.balance)
                                        .split(".")
                                        .collect::<Vec<&str>>()[1]
                                )
                            } else {
                                ".??".to_string()
                            }
                        })
                        .size(20),
                    ]
                    .padding(0)
                    .align_items(Alignment::End)
                ),
                Space::with_height(10),
                button("Sync").on_press(Message::Sync),
                Space::with_height(20),
                rule::Rule::horizontal(0.0),
                Space::with_height(20),
                text("Send").size(23),
                Space::with_height(10),
                text("Amount").size(15),
                text_input("Enter an amount...", self.amount.to_string().as_str())
                    .on_input(Message::AmountChanged)
                    .padding(10)
                    .size(20),
                Space::with_height(10),
                text("Address").size(15),
                text_input("Enter an address...", self.address.as_str())
                    .on_input(Message::AddressChanged)
                    .padding(10)
                    .size(20),
                Space::with_height(15),
                button("Send").on_press(Message::Send),
                text("Public Key:").size(20),
                container(
                    row![
                        text(crate::shorten::shorten(crate::key::get_public_key().expect("Failed to get public key"))).size(15),
                        button("Copy").on_press(Message::Copy)
                    ]
                )
            ]
            .padding(20)
            .align_items(Alignment::Start),
        )
        .into()
    }

    fn theme(&self) -> iced::Theme {
        iced::Theme::Dark
    }
}
