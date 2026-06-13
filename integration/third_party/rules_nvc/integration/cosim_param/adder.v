module adder #(
    parameter WIDTH = 8
)(
    input  wire clk,
    input  wire [WIDTH-1:0] a,
    input  wire [WIDTH-1:0] b,
    output reg  [WIDTH:0] sum
);

    always @(posedge clk) begin
        sum <= {1'b0, a} + {1'b0, b};
    end

endmodule
