use iced::widget::{button, row, column, text, container, text_input};
use iced::{Alignment, Element, Sandbox, Settings};
mod sync;

pub fn main() -> iced::Result {
    App::run(Settings::default())
}

struct App {
    balance: f64,
    amount: f64,
}

#[derive(Debug, Clone)]
enum Message {
    Sync,
    AmountChanged(String),
}

impl Sandbox for App {
    type Message = Message;

    fn new() -> Self {
        Self {
            balance: 0.0,
            amount: 0.0,
        }
    }

    fn title(&self) -> String {
        String::from("GUI Wallet")
    }

    fn update(&mut self, message: Message) {
        match message {
            Message::Sync => {
                self.balance = crate::sync::sync();
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
        }
    }

    fn view(&self) -> Element<Message> {
        column![
            container(
                row![
                    text(format!("{:.0}", self.balance)).size(50),
                    text(format!(".{}", format!("{:.2}", self.balance).split(".").collect::<Vec<&str>>()[1])).size(20),
                ]
                .padding(20)
                .align_items(Alignment::End)
            ),
            button("Sync").on_press(Message::Sync),
            text_input("Enter an amount...", self.amount.to_string().as_str())
                .on_input(Message::AmountChanged)
                .padding(10)
                .size(30),
        ]
        .padding(20)
        .align_items(Alignment::Center)
        .into()
    }
    
    fn theme(&self) -> iced::Theme {
        iced::Theme::Dark
    }
}