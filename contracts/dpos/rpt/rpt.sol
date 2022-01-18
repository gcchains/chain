

pragma solidity ^0.4.24;


contract Rpt {
    // The 5 weight configs.
    uint public alpha = 50;
    uint public beta = 15;
    uint public gamma = 10;
    uint public psi = 15;
    uint public omega = 10;
    
    // other configs.
    uint public window = 100; // number of blocks used for rpt calculation

    // 3 election configs
    uint public lowRptPercentage = 70; // percentage: 0-100
    uint public totalSeats = 8; // 0-8
    uint public lowRptSeats = 2; // 0-8 && lower than totalSeats

    address owner;
    
    modifier onlyOwner() {require(msg.sender == owner);_;}
    
    event UpdateWeightConfigs(uint blockNumber);
    event UpdateOneConfig(uint blockNumber, string configName, uint configValue);
    event UpdateElectionConfigs(uint blockNumber);
    
    
    constructor() public {
        owner = msg.sender;
    }

    
    /** modified all configs. */
    function updateWeightConfigs(uint _alpha, uint _beta, uint _gamma, uint _psi, uint _omega, uint _window)
        public 
        onlyOwner()
        {
            require(_window>=10 && _window<=100);
            require(_alpha<=100 && _beta<=100 && _gamma<=100 && _psi<=100 && _omega<=100);
            alpha = _alpha;
            beta = _beta;
            gamma = _gamma;
            psi = _psi;
            omega = _omega;
            window = _window;

            emit UpdateWeightConfigs(block.number);
        }


    function updateElectionConfigs(uint _lowRptPercentage, uint _totalSeats, uint _lowRptSeats) public onlyOwner {
        require(_lowRptPercentage<=100 && _totalSeats<=8 && _lowRptSeats<=_totalSeats);
        lowRptPercentage = _lowRptPercentage;
        totalSeats = _totalSeats;
        lowRptSeats = _lowRptSeats;

        emit UpdateElectionConfigs(block.number);
    }

    function updateLowRptPercentage(uint _lowRptPercentage) public onlyOwner {
        require(_lowRptPercentage <= 100);
        lowRptPercentage = _lowRptPercentage;
        emit UpdateOneConfig(block.number, "lowRptPercentage", lowRptPercentage);
    }

    function updateTotalSeats(uint _totalSeats) public onlyOwner {
        require(_totalSeats<=8 && _totalSeats>=lowRptSeats);
        totalSeats = _totalSeats;
        emit UpdateOneConfig(block.number, "totalSeats", totalSeats);
    }

    function updateLowRptSeats(uint _lowRptSeats) public onlyOwner {
        require(_lowRptSeats <= totalSeats);
        lowRptSeats = _lowRptSeats;
        emit UpdateOneConfig(block.number, "lowRptSeats", lowRptSeats);
    }
    
    /** modified one config. */
    function updateAlpha(uint _alpha) public onlyOwner {
        require(_alpha <= 100);
        alpha = _alpha;
        emit UpdateOneConfig(block.number, "alpha", alpha);
    }
    
    function updateBeta(uint _beta) public onlyOwner {
        require(_beta <= 100);
        beta = _beta;
        emit UpdateOneConfig(block.number, "beta", beta);
    }
    
    function updateGamma(uint _gamma) public onlyOwner {
        require(_gamma <= 100);
        gamma = _gamma;
        emit UpdateOneConfig(block.number, "gamma", gamma);
    }
    
    function updatePsi(uint _psi) public onlyOwner {
        require(_psi <= 100);
        psi = _psi;
        emit UpdateOneConfig(block.number, "psi", psi);
    }
    
    function updateOmega(uint _omega) public onlyOwner {
        require(_omega <= 100);
        omega = _omega;
        emit UpdateOneConfig(block.number, "omega", omega);
    }
    
    function updateWindow(uint _window) public onlyOwner {
        require(_window>=10 && _window<=100);
        window = _window;
        emit UpdateOneConfig(block.number, "window", window);
    }

}