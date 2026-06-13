/// This is a basic verilog counter, for testing NVC verilog support.
module counter #(
    parameter int Width = 8
) (
    input logic reset, ///< Synchronous reset, active high.
    input logic clk, ///< Clock.
    output reg [Width-1:0] count ///< The counter output value.
);
    logic [Width-1:0] next_value;
    always_comb begin
        next_value = count + 1;
    end
    always_ff @(posedge clk) begin
        if (reset == 1'b1) begin
            count <= 1'b0;
        end else begin
            count <= next_value;
        end
    end
endmodule
