use iced::widget::{button, row, column, text, container, text_input, rule, Space};
use iced::{Alignment, Element, Sandbox, Settings};
use crate::send::send;

mod sync;
mod send;

pub fn main() -> iced::Result {
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
        }
    }

    fn view(&self) -> Element<Message> {
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
                    }).size(50),
                    text({
                        if self.balance != -1.0 {
                            format!(".{}", format!("{:.2}", self.balance).split(".").collect::<Vec<&str>>()[1])
                        } else {
                            ".??".to_string()
                        }
                    }).size(20),
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
            button("Send").on_press(Message::Send)
        ]
        .padding(20)
        .align_items(Alignment::Start)
        .into()
    }
    
    fn theme(&self) -> iced::Theme {
        iced::Theme::Dark
    }
}